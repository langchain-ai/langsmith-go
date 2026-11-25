package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/spf13/cobra"
)

// wholeFileOnlyFlags are the upload flags that --function does not accept.
var wholeFileOnlyFlags = []string{"advanced", "dependencies", "workspace-secret", "input", "no-wait", "ignore-import-check"}

// How often and how long upload polls an advanced evaluator's build.
var (
	evaluatorBuildPollInterval = 2 * time.Second
	evaluatorBuildTimeout      = 15 * time.Minute
)

// wholeFileUpload is an upload without --function: the file is sent unchanged
// through the evaluators API and the rule attaches the evaluator by ID.
type wholeFileUpload struct {
	file, name        string
	target            evaluatorTarget
	samplingRate      float64
	traceFilter       string
	replace, yes      bool
	advanced, noWait  bool
	dependenciesFile  string
	ignoreCredentials bool
	ignoreImports     bool
	workspaceSecrets  []string
	input             string
}

// checkFunctionFlagCombination rejects the whole-file flags alongside --function.
func checkFunctionFlagCombination(cmd *cobra.Command) error {
	for _, name := range wholeFileOnlyFlags {
		if cmd.Flags().Changed(name) {
			return fmt.Errorf("--%s cannot be combined with --function; omit --function to upload the whole file", name)
		}
	}
	return nil
}

func uploadWholeFile(ctx context.Context, cmd *cobra.Command, c *client.Client, u wholeFileUpload) error {
	source, err := os.ReadFile(u.file)
	if err != nil {
		return fmt.Errorf("reading evaluator file: %w", err)
	}
	language, _ := detectLanguage(u.file)
	if language == "" {
		return fmt.Errorf("unsupported file extension: %s (use .py, .js, .ts, .tsx, or .mjs)", u.file)
	}
	var dependencies string
	if u.dependenciesFile != "" {
		b, err := os.ReadFile(u.dependenciesFile)
		if err != nil {
			return fmt.Errorf("reading dependencies file: %w", err)
		}
		dependencies = string(b)
		if err := checkForSecrets(u.dependenciesFile, dependencies, u.ignoreCredentials); err != nil {
			return err
		}
	}
	if err := checkForSecrets(u.file, string(source), u.ignoreCredentials); err != nil {
		return err
	}

	existing, err := findRuleToReplace(ctx, cmd, c, u.name, u.target, u.replace, u.yes)
	if err != nil {
		return err
	}
	// The server checks a legacy evaluator's imports when it saves the code,
	// but an advanced evaluator's build only checks syntax, so a missing local
	// module would surface only when the evaluator runs.
	advanced := u.advanced
	if !advanced && existing != nil && existing.EvaluatorID != "" {
		if advanced, err = isAdvancedEvaluator(ctx, c, existing.EvaluatorID); err != nil {
			return err
		}
	}
	if advanced {
		if err := checkLocalImports(u.file, string(source), language); err != nil {
			if !u.ignoreImports {
				return fmt.Errorf("%w\nIf these modules will be available when the evaluator runs, pass --ignore-import-check", err)
			}
			fmt.Fprintf(os.Stderr, "warning: %v (--ignore-import-check)\n", err)
		}
	}

	var ev langsmith.OnlineEvaluator
	var ruleID string
	if existing == nil {
		ev, ruleID, err = createAndAttachEvaluator(ctx, c, u, string(source), language, dependencies)
	} else {
		ev, err = replaceEvaluatorCode(ctx, cmd, c, u, *existing, string(source), language, dependencies)
		ruleID = existing.ID
	}
	if err != nil {
		return err
	}

	result := map[string]any{
		"status":       "uploaded",
		"rule_id":      ruleID,
		"evaluator_id": ev.ID,
		"name":         u.name,
		"target":       targetKind(u.target),
		"advanced":     ev.CodeEvaluator.AdvancedFeaturesEnabled,
	}
	if ev.CodeEvaluator.AdvancedFeaturesEnabled {
		built := &ev
		if !u.noWait && buildInProgress(ev.CodeEvaluator.EvaluatorBuildStatus) {
			if built, err = waitForEvaluatorBuild(ctx, c, ev.ID); err != nil {
				return err
			}
		}
		// An update that changes neither code nor dependencies starts no build,
		// so the evaluator can come back already FAILED.
		status := built.CodeEvaluator.EvaluatorBuildStatus
		if status == langsmith.OnlineCodeEvaluatorEvaluatorBuildStatusFailed {
			return fmt.Errorf("evaluator %s build failed: %s\nThe evaluator and rule %s were saved; fix the code or dependencies and rerun with --replace", ev.ID, built.CodeEvaluator.EvaluatorBuildError, ruleID)
		}
		result["build_status"] = string(status)
	}
	return output.OutputJSON(result, "")
}

func codeEvaluatorParams(u wholeFileUpload, code, language, dependencies string) langsmith.CreateOnlineCodeEvaluatorRequestParam {
	p := langsmith.CreateOnlineCodeEvaluatorRequestParam{
		Code:     langsmith.F(code),
		Language: langsmith.F(language),
	}
	// advanced_features_enabled is sent only with --advanced: the server rejects
	// the field, even false, on workspaces without code evaluators v2.
	if u.advanced {
		p.AdvancedFeaturesEnabled = langsmith.F(true)
	}
	if u.dependenciesFile != "" {
		p.Dependencies = langsmith.F(dependencies)
	}
	if len(u.workspaceSecrets) > 0 {
		p.WorkspaceSecretsKeys = langsmith.F(u.workspaceSecrets)
	}
	if u.input != "" {
		p.CodeEvaluatorInput = langsmith.F(langsmith.CreateOnlineCodeEvaluatorRequestCodeEvaluatorInput(u.input))
	}
	return p
}

// createAndAttachEvaluator creates the evaluator, then a rule that attaches it
// to the target. The evaluator is deleted again when the rule cannot be created.
func createAndAttachEvaluator(ctx context.Context, c *client.Client, u wholeFileUpload, code, language, dependencies string) (langsmith.OnlineEvaluator, string, error) {
	resp, err := c.SDK.OnlineEvaluators.New(ctx, langsmith.OnlineEvaluatorNewParams{
		CreateOnlineEvaluatorRequest: langsmith.CreateOnlineEvaluatorRequestParam{
			Name:          langsmith.F(u.name),
			Type:          langsmith.F(langsmith.OnlineEvaluatorTypeCode),
			CodeEvaluator: langsmith.F(codeEvaluatorParams(u, code, language, dependencies)),
		},
	})
	if err != nil {
		return langsmith.OnlineEvaluator{}, "", fmt.Errorf("creating evaluator: %w", withEntryPointHint(err, code, language, u.input))
	}
	ev := resp.Evaluator

	payload := map[string]any{
		"display_name":           u.name,
		"sampling_rate":          u.samplingRate,
		"is_enabled":             true,
		"include_extended_stats": false,
		"evaluator_id":           ev.ID,
		"evaluator_version":      3,
	}
	addRuleTarget(payload, u.target)
	if u.traceFilter != "" {
		payload["trace_filter"] = u.traceFilter
	}
	if groupBy := groupByForInput(u.input); groupBy != "" {
		payload["group_by"] = groupBy
	}
	var result map[string]any
	if err := c.RawPost(ctx, "/api/v1/runs/rules", payload, &result); err != nil {
		attachErr := fmt.Errorf("attaching evaluator to the %s: %w", targetKind(u.target), err)
		if delErr := deleteEvaluator(ctx, c, ev.ID, !u.noWait); delErr != nil {
			return ev, "", fmt.Errorf("%w\nevaluator %s was created but could not be deleted: %v", attachErr, ev.ID, delErr)
		}
		return ev, "", fmt.Errorf("%w\nevaluator %s was deleted", attachErr, ev.ID)
	}
	ruleID, _ := result["id"].(string)
	return ev, ruleID, nil
}

// replaceEvaluatorCode updates the code of the evaluator behind existing and
// sends the rule's settings back, as PATCH /runs/rules/{id} is a full replace.
// The evaluator keeps its mode unless --advanced is set.
func replaceEvaluatorCode(ctx context.Context, cmd *cobra.Command, c *client.Client, u wholeFileUpload, existing langsmith.Rule, code, language, dependencies string) (langsmith.OnlineEvaluator, error) {
	if existing.EvaluatorID == "" {
		return langsmith.OnlineEvaluator{}, fmt.Errorf("rule %q (%s) has no evaluator ID, so its code cannot be replaced through the evaluators API; delete the rule with 'langsmith evaluator rule delete' and upload again", existing.DisplayName, existing.ID)
	}
	// The rules API rejects moving a rule into or out of trajectory grouping.
	// The SDK's RuleGroupBy constants are upper case; the API returns lower case.
	groupBy := groupByForInput(u.input)
	if u.input != "" && (groupBy == "trajectory_id") != strings.EqualFold(string(existing.GroupBy), "trajectory_id") {
		return langsmith.OnlineEvaluator{}, fmt.Errorf("rule %q (%s) cannot switch to or from trajectory input; delete the rule with 'langsmith evaluator rule delete' and upload again", existing.DisplayName, existing.ID)
	}
	// The rule payload is built before any write so that a rule the CLI
	// refuses to update leaves the evaluator unchanged too.
	payload := map[string]any{"sampling_rate": u.samplingRate}
	addRuleTarget(payload, u.target)
	if u.traceFilter != "" {
		payload["trace_filter"] = u.traceFilter
	}
	if groupBy != "" {
		payload["group_by"] = groupBy
	}
	if err := keepRuleSettings(payload, existing, replaceOverridesFromFlags(cmd)); err != nil {
		return langsmith.OnlineEvaluator{}, err
	}

	create := codeEvaluatorParams(u, code, language, dependencies)
	update := langsmith.UpdateOnlineCodeEvaluatorRequestParam{
		Code:                    create.Code,
		Language:                create.Language,
		AdvancedFeaturesEnabled: create.AdvancedFeaturesEnabled,
		Dependencies:            create.Dependencies,
		WorkspaceSecretsKeys:    create.WorkspaceSecretsKeys,
	}
	if u.input != "" {
		update.CodeEvaluatorInput = langsmith.F(langsmith.UpdateOnlineCodeEvaluatorRequestCodeEvaluatorInput(u.input))
	}
	params := langsmith.OnlineEvaluatorUpdateParams{
		UpdateOnlineEvaluatorRequest: langsmith.UpdateOnlineEvaluatorRequestParam{
			CodeEvaluator: langsmith.F(update),
		},
	}
	resp, err := c.SDK.OnlineEvaluators.Update(ctx, existing.EvaluatorID, params)
	if isBuildInProgress(err) && !u.noWait {
		if _, waitErr := waitForEvaluatorBuild(ctx, c, existing.EvaluatorID); waitErr != nil {
			return langsmith.OnlineEvaluator{}, waitErr
		}
		resp, err = c.SDK.OnlineEvaluators.Update(ctx, existing.EvaluatorID, params)
	}
	if err != nil {
		return langsmith.OnlineEvaluator{}, fmt.Errorf("replacing evaluator %s: %w", existing.EvaluatorID, withEntryPointHint(err, code, language, u.input))
	}

	var result map[string]any
	if err := c.RawPatch(ctx, fmt.Sprintf("/api/v1/runs/rules/%s", existing.ID), payload, &result); err != nil {
		return resp.Evaluator, fmt.Errorf("evaluator %s was updated, but updating rule %s failed: %w", existing.EvaluatorID, existing.ID, err)
	}
	return resp.Evaluator, nil
}

// refuseFunctionReplaceOfAdvanced stops --replace --function from overwriting an
// advanced evaluator's whole file with one function body.
func refuseFunctionReplaceOfAdvanced(ctx context.Context, c *client.Client, existing langsmith.Rule) error {
	if existing.EvaluatorID == "" {
		return nil
	}
	advanced, err := isAdvancedEvaluator(ctx, c, existing.EvaluatorID)
	if err != nil {
		return err
	}
	if advanced {
		return fmt.Errorf("evaluator %q is an advanced evaluator; --function would replace its whole file with one function, so upload the file without --function", existing.EvaluatorName)
	}
	return nil
}

func isAdvancedEvaluator(ctx context.Context, c *client.Client, id string) (bool, error) {
	ev, err := c.SDK.OnlineEvaluators.Get(ctx, id)
	if err != nil {
		return false, fmt.Errorf("fetching evaluator %s: %w", id, err)
	}
	return ev.CodeEvaluator.AdvancedFeaturesEnabled, nil
}

// deleteEvaluator removes an evaluator. An advanced evaluator cannot be deleted
// while it builds, so with wait set it waits for the build and tries again.
func deleteEvaluator(ctx context.Context, c *client.Client, id string, wait bool) error {
	err := c.SDK.OnlineEvaluators.Delete(ctx, id, langsmith.OnlineEvaluatorDeleteParams{})
	if isBuildInProgress(err) && wait {
		if _, waitErr := waitForEvaluatorBuild(ctx, c, id); waitErr != nil {
			return waitErr
		}
		err = c.SDK.OnlineEvaluators.Delete(ctx, id, langsmith.OnlineEvaluatorDeleteParams{})
	}
	return err
}

// waitForEvaluatorBuild polls the evaluator until its build is READY or FAILED.
func waitForEvaluatorBuild(ctx context.Context, c *client.Client, id string) (*langsmith.OnlineEvaluator, error) {
	deadline := time.Now().Add(evaluatorBuildTimeout)
	for announced := false; ; announced = true {
		ev, err := c.SDK.OnlineEvaluators.Get(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("checking build of evaluator %s: %w", id, err)
		}
		if !buildInProgress(ev.CodeEvaluator.EvaluatorBuildStatus) {
			return ev, nil
		}
		if !announced {
			fmt.Fprintf(os.Stderr, "Waiting for evaluator %s to build (--no-wait skips this)...\n", id)
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("evaluator %s is still %s after %s; check it later with 'langsmith evaluator get %s'", id, ev.CodeEvaluator.EvaluatorBuildStatus, evaluatorBuildTimeout, id)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(evaluatorBuildPollInterval):
		}
	}
}

func buildInProgress(status langsmith.OnlineCodeEvaluatorEvaluatorBuildStatus) bool {
	return status == langsmith.OnlineCodeEvaluatorEvaluatorBuildStatusEnqueued ||
		status == langsmith.OnlineCodeEvaluatorEvaluatorBuildStatusBuilding
}

// isBuildInProgress reports whether err is the 409 the evaluators API returns
// while a build is enqueued or running.
func isBuildInProgress(err error) bool {
	var apiErr *langsmith.Error
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusConflict
}

// groupByForInput is the rule grouping that delivers an input mode's data:
// trajectory input runs on trajectories, the other modes on threads.
func groupByForInput(input string) string {
	switch input {
	case "":
		return ""
	case string(langsmith.CreateOnlineCodeEvaluatorRequestCodeEvaluatorInputTrajectory):
		return "trajectory_id"
	default:
		return "thread_id"
	}
}

func addRuleTarget(payload map[string]any, target evaluatorTarget) {
	if target.datasetID != "" {
		payload["dataset_id"] = target.datasetID
	}
	if target.projectID != "" {
		payload["session_id"] = target.projectID
	}
}

func targetKind(target evaluatorTarget) string {
	if target.datasetID != "" {
		return "dataset"
	}
	return "project"
}
