package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/langchain-ai/langsmith-go/cli/internal/cache"
	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	"github.com/langchain-ai/langsmith-go/cli/internal/cmd/api"
	"github.com/langchain-ai/langsmith-go/cli/internal/cmdutil"
	lsconfig "github.com/langchain-ai/langsmith-go/cli/internal/config"
	"github.com/langchain-ai/langsmith-go/cli/internal/exitcode"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/langchain-ai/langsmith-go/cli/internal/usage"
	"github.com/spf13/cobra"
)

// Global flag values stored here; accessed by subcommands via helpers.
var (
	flagAPIKey       string
	flagAPIURL       string
	flagProfile      string
	flagWorkspaceID  string
	flagOutputFormat string
	flagJSON         bool
)

const feedbackHintInterval = 24 * time.Hour

// NewRootCmd creates the top-level `langsmith` command.
func NewRootCmd(rawVersion, displayVersion string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "langsmith",
		Short: "LangSmith CLI — query and manage LangSmith resources",
		Long: `LangSmith CLI — query and manage LangSmith resources from the command line.

	Designed for developers and AI coding agents who need fast access to traces,
	runs, datasets, evaluators, experiments, and threads.

Authentication:
  Run 'langsmith auth login' or set LANGSMITH_API_KEY.
  Optionally set LANGSMITH_ENDPOINT for self-hosted instances.
  Use --profile or LANGSMITH_PROFILE to select a saved profile.
  Pass --workspace to target a specific workspace for one command.
  Set a default workspace with 'langsmith profiles set-workspace <workspace-id>'.
  Set LANGSMITH_PROJECT as a default project name for trace/run queries.

Quick start:
  langsmith projects list
  langsmith traces list --project my-project --limit 5
  langsmith runs list --project my-project --run-type llm --limit 10
  langsmith datasets list
  langsmith evaluators list
  langsmith experiments list --dataset my-eval-dataset

Common workflows:
  Find and confirm a fix for production failures:
    langsmith investigate --project my-app --last-n-minutes 60
    langsmith verify --project my-app --since 2026-11-10T14:00:00Z
  Turn failed runs into an evaluation:
    langsmith runs list --project my-app --error --include-io --limit 20
    langsmith datasets create --name my-eval-dataset
    langsmith examples create --dataset my-eval-dataset --inputs '{"question": "..."}' --outputs '{"answer": "..."}'
    langsmith evaluators upload evals.py --name correctness --dataset my-eval-dataset
    langsmith experiments list --dataset my-eval-dataset
    langsmith experiments get my-experiment

Coding agents: run 'langsmith docs agents' for task-by-task usage.

	Output:
	  --format pretty  Human-readable tables, trees, and syntax-highlighted JSON (default).
	  --format json    Machine-readable JSON for agents and scripts. From 1.0.0 it is
	                   the default when stdout is piped or redirected.
	  --json           Shorthand for --format json.
	  --jq EXPR        Filter JSON output with a jq expression (implies --format json).

Exit codes (every failure exits 1 until 1.0.0; set LANGSMITH_EXIT_CODES=detailed
to use these now):
  0  success
  1  bad input (fix it before retrying), or a failure not listed below
  2  missing or rejected credentials (HTTP 401, 403)
  3  not found (HTTP 404, 410)
  4  conflict or plan limit (HTTP 409, 402); don't retry unchanged
  5  a --wait timed out; the operation may still finish, so check its
     status instead of starting it again
  6  server error or timeout (HTTP 5xx, 408), or the API could not be
     reached; retry reads with backoff, and check whether a write took
     effect before repeating it
  7  rate limited (HTTP 429); retry with backoff`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       displayVersion,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			usage.SetCommand(cmd)
			cmdutil.ResetFormatWarning(cmd.ErrOrStderr())
			cmdutil.ResetNotices(cmd.ErrOrStderr())
			if err := cmdutil.NormalizeFormatFlags(cmd); err != nil {
				return err
			}
			// Reject a bad --jq before the command does any work.
			expr := cmdutil.ResolveJQ(cmd)
			if expr != "" {
				if _, err := output.CompileJQ(expr); err != nil {
					return err
				}
			}
			output.SetJQ(expr)
			warnDeprecated(cmd.ErrOrStderr(), cmd, args, os.Args[1:], output.JQ() != "" || flagOutputFormat == "json")
			return nil
		},
	}
	usage.SetVersion(rawVersion)

	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "LangSmith API key [env: LANGSMITH_API_KEY]")
	_ = rootCmd.PersistentFlags().MarkHidden("api-key")
	rootCmd.PersistentFlags().StringVar(&flagAPIURL, "api-url", "", "LangSmith API URL [env: LANGSMITH_ENDPOINT]")
	rootCmd.PersistentFlags().StringVar(&flagProfile, "profile", "", "Named profile to use [env: LANGSMITH_PROFILE]")
	rootCmd.PersistentFlags().StringVar(&flagWorkspaceID, "workspace", "", "LangSmith workspace ID [env: LANGSMITH_WORKSPACE_ID]")
	rootCmd.PersistentFlags().StringVar(&flagWorkspaceID, "workspace-id", "", "LangSmith workspace ID [env: LANGSMITH_WORKSPACE_ID]")
	_ = rootCmd.PersistentFlags().MarkHidden("workspace-id")
	rootCmd.PersistentFlags().StringVar(&flagOutputFormat, "format", "", "Output format: pretty or json (default pretty; json when piped from 1.0.0)")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, cmdutil.JSONFlagUsage)
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		showFeedbackHint(cmd)
		return err
	})

	// Register all subcommand groups
	rootCmd.AddCommand(newProductFeedbackCmd(rawVersion))
	rootCmd.AddCommand(newProjectCmd())
	rootCmd.AddCommand(newTraceCmd())
	rootCmd.AddCommand(newAgentSetupCmd())
	rootCmd.AddCommand(newRunCmd())
	rootCmd.AddCommand(newChartCmd())
	rootCmd.AddCommand(newThreadCmd())
	rootCmd.AddCommand(newDatasetCmd())
	rootCmd.AddCommand(newExampleCmd())
	rootCmd.AddCommand(newEvaluatorCmd())
	rootCmd.AddCommand(newRulesAliasCmd())
	rootCmd.AddCommand(newExperimentCmd())
	rootCmd.AddCommand(newSandboxCmd())
	rootCmd.AddCommand(newTunnelCmd())
	rootCmd.AddCommand(newRelayCmd())
	rootCmd.AddCommand(newConnectorCmd(rawVersion, displayVersion))
	rootCmd.AddCommand(newInsightsCmd())
	rootCmd.AddCommand(newInvestigateCmd())
	rootCmd.AddCommand(newVerifyCmd())
	rootCmd.AddCommand(newHubCmd())
	rootCmd.AddCommand(newAppsCmd())
	rootCmd.AddCommand(newPromptCmd())
	rootCmd.AddCommand(authCommand.Cobra())
	rootCmd.AddCommand(gatewayCommand.Cobra())
	rootCmd.AddCommand(newProfileCmd())
	rootCmd.AddCommand(newWorkspaceCmd())
	rootCmd.AddCommand(newUpdateCmd(rawVersion))
	rootCmd.AddCommand(newDocsCmd())
	rootCmd.AddCommand(api.NewCmd())
	addGeneratedCommands(rootCmd)
	// Cobra would add the completion group at Execute time; add it now so the
	// walk below covers it too. Its commands print shell scripts, not JSON, so
	// they take no --jq.
	rootCmd.InitDefaultCompletionCmd()
	if completion, _, err := rootCmd.Find([]string{"completion"}); err == nil && completion != rootCmd {
		for _, shell := range completion.Commands() {
			if shell.Annotations == nil {
				shell.Annotations = map[string]string{}
			}
			shell.Annotations[cmdutil.NoJSONOutputAnnotation] = "true"
		}
	}
	rejectUnknownSubcommands(rootCmd)
	cmdutil.AddJQFlags(rootCmd)

	return rootCmd
}

// rejectUnknownSubcommands makes every command group without an action of its
// own fail on an unknown subcommand, as the root does. Cobra checks subcommand
// names only at the root; a nested group without Run prints its help and exits
// 0, which hides a mistyped or invented command. A bare group still prints help.
func rejectUnknownSubcommands(cmd *cobra.Command) {
	if cmd.HasParent() && cmd.HasSubCommands() && !cmd.Runnable() {
		cmd.Args = unknownSubcommandArgs
		cmd.RunE = func(cmd *cobra.Command, _ []string) error { return cmd.Help() }
		if cmd.Annotations == nil {
			cmd.Annotations = map[string]string{}
		}
		cmd.Annotations[cmdutil.NoJSONOutputAnnotation] = "true"
	}
	for _, c := range cmd.Commands() {
		rejectUnknownSubcommands(c)
	}
}

// unknownSubcommandArgs rejects any positional argument to a command group,
// with cobra's root-level message and "Did you mean" suggestions.
func unknownSubcommandArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	if cmd.SuggestionsMinimumDistance <= 0 {
		cmd.SuggestionsMinimumDistance = 2
	}
	msg := fmt.Sprintf("unknown command %q for %q", args[0], cmd.CommandPath())
	if suggestions := cmd.SuggestionsFor(args[0]); len(suggestions) > 0 && !cmd.DisableSuggestions {
		msg += "\n\nDid you mean this?\n\t" + strings.Join(suggestions, "\n\t")
	}
	return errors.New(msg)
}

// GetAPIKey resolves the API key from flag → env → profile.
func GetAPIKey() string {
	opts, _ := resolveClientOptions()
	return opts.APIKey
}

// apiKeyRefProfile returns the selected profile's name when the active
// credentials are its api_key_env or api_key_cmd, so a caller can hand the
// profile on rather than the resolved key.
func apiKeyRefProfile() (string, bool) {
	if flagAPIKey != "" || os.Getenv("LANGSMITH_API_KEY") != "" {
		return "", false
	}
	cfg, err := lsconfig.Load()
	if err != nil {
		return "", false
	}
	name, p, ok := cfg.ResolveProfile(flagProfile, profileEnvName())
	if !ok || p.AccessToken() != "" || p.OAuth.RefreshToken != "" || p.APIKey != "" {
		return "", false
	}
	return name, p.APIKeyEnv != "" || p.APIKeyCmd != ""
}

// GetAPIURL resolves the API URL from flag → env → profile → default.
func GetAPIURL() string {
	opts, _ := resolveClientOptions()
	return opts.APIURL
}

// GetWorkspaceID resolves the workspace ID from flag → env → profile.
func GetWorkspaceID() string {
	opts, _ := resolveClientOptions()
	return opts.WorkspaceID
}

// GetFormat returns the output format. --jq filters JSON, so it implies json;
// without either flag, cmdutil.DefaultFormat applies.
func GetFormat() string {
	return resolveFormat(cmdutil.DefaultFormat)
}

// resolveFormat is GetFormat with defaultFormat picking the format when no
// flag does.
func resolveFormat(defaultFormat func() string) string {
	if output.JQ() != "" {
		return "json"
	}
	if flagOutputFormat != "" {
		return flagOutputFormat
	}
	return defaultFormat()
}

func showFeedbackHint(cmd *cobra.Command) {
	// Flag errors arrive before --format is resolved, so only an explicit
	// non-pretty format suppresses the hint.
	if (flagOutputFormat != "" && flagOutputFormat != "pretty") || flagJSON || cmd == nil {
		return
	}
	path := cmd.CommandPath()
	if path == "langsmith" || path == "langsmith report" || path == "langsmith feedback" || path == "langsmith auth" || path == "langsmith auth login" {
		return
	}
	marker := filepath.Join(cache.DefaultDir(), "feedback-hint")
	if info, err := os.Stat(marker); err == nil && time.Since(info.ModTime()) < feedbackHintInterval {
		return
	}
	fmt.Fprintln(cmd.ErrOrStderr(), "Tip: report CLI friction with 'langsmith report \"note\"'.")
	_ = cache.Write(marker, nil)
}

// MustGetClient creates a LangSmith client or exits with an error.
func MustGetClient() *client.Client {
	c, err := getClient()
	if err != nil {
		ExitErrorf("%v", err)
	}
	return c
}

// getClient is the non-exiting sibling of MustGetClient: it returns an error
// instead of calling os.Exit, so request handlers can't crash.
func getClient() (*client.Client, error) {
	opts, err := resolveAuthenticatedOptions(currentGlobalFlags())
	if err != nil {
		return nil, err
	}
	cmdutil.NoteActiveProfile(flagProfile, opts, noticesAsJSON())
	return client.NewWithOptions(opts), nil
}

// noticesAsJSON reports whether stderr notices are JSON objects: they are
// whenever the output is, so they follow GetFormat and its default for
// output that is not a terminal. It doesn't print the default-format warning.
func noticesAsJSON() bool {
	return strings.EqualFold(resolveFormat(cmdutil.DefaultFormatQuiet), "json")
}

// globalFlags holds the values of langsmith's global auth and routing flags.
type globalFlags struct {
	APIKey      string
	APIURL      string
	Profile     string
	WorkspaceID string
}

// currentGlobalFlags returns the global flags Cobra parsed.
func currentGlobalFlags() globalFlags {
	return globalFlags{APIKey: flagAPIKey, APIURL: flagAPIURL, Profile: flagProfile, WorkspaceID: flagWorkspaceID}
}

// selectionArgs returns the flags that chose the profile, endpoint and
// workspace, formatted to append to a suggested follow-up command so it acts
// on the same workspace. The API key is never included.
func (f globalFlags) selectionArgs() string {
	var b strings.Builder
	for _, kv := range [][2]string{{"--profile", f.Profile}, {"--api-url", f.APIURL}, {"--workspace", f.WorkspaceID}} {
		if kv[1] != "" {
			fmt.Fprintf(&b, " %s %s", kv[0], shellQuoteIfNeeded(kv[1]))
		}
	}
	return b.String()
}

// shellQuoteIfNeeded leaves plain values bare so suggested commands stay
// readable, and single-quotes anything the shell would split or expand.
func shellQuoteIfNeeded(s string) string {
	plain := strings.IndexFunc(s, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("-_.:/@=+,%", r))
	}) == -1
	if plain {
		return s
	}
	return shellQuote(s)
}

// resolveAuthenticatedOptions resolves the client options and fails when they
// carry no credentials.
func resolveAuthenticatedOptions(flags globalFlags) (client.Options, error) {
	opts, err := resolveClientOptionsFor(flags)
	if err != nil {
		return opts, err
	}
	if !opts.HasAuth() {
		return opts, exitcode.WithCode(fmt.Errorf("not authenticated; run 'langsmith auth login', set LANGSMITH_API_KEY, or pass --api-key"), exitcode.Auth)
	}
	return opts, nil
}

func resolveClientOptions() (client.Options, error) {
	return resolveClientOptionsFor(currentGlobalFlags())
}

// resolveClientOptionsFor resolves auth, endpoint, and workspace from the given
// global flags, then the environment, then the selected profile.
func resolveClientOptionsFor(flags globalFlags) (client.Options, error) {
	opts := client.Options{APIURL: lsconfig.DefaultAPIURL}

	cfg, err := lsconfig.Load()
	var cfgErr error
	if err != nil {
		cfgErr = err
		cfg = &lsconfig.Config{Profiles: make(map[string]lsconfig.Profile)}
	}

	envProfile := strings.TrimSpace(os.Getenv("LANGSMITH_PROFILE"))
	profileName, profile, hasProfile := "", lsconfig.Profile{}, false
	if flags.Profile != "" || envProfile != "" || cfgErr == nil {
		if cfgErr != nil && (flags.Profile != "" || envProfile != "") {
			return opts, cfgErr
		}
		profileName, profile, hasProfile = cfg.ResolveProfile(flags.Profile, envProfile)
		if (flags.Profile != "" || envProfile != "") && !hasProfile {
			return opts, fmt.Errorf("profile not found: %s (run `langsmith profiles list`)", profileName)
		}
	}

	if hasProfile {
		if profile.APIURL != "" {
			opts.APIURL = profile.APIURL
		}
		opts.WorkspaceID = profile.WorkspaceID
	}

	if v := os.Getenv("LANGSMITH_ENDPOINT"); v != "" {
		if flags.Profile != "" && hasProfile && profile.APIURL != "" {
			fmt.Fprintf(os.Stderr, "warning: ignoring LANGSMITH_ENDPOINT because profile %q was selected with --profile\n", profileName)
		} else {
			opts.APIURL = client.NormalizeURL(v)
		}
	}
	if flags.APIURL != "" {
		opts.APIURL = client.NormalizeURL(flags.APIURL)
	}

	if v := os.Getenv("LANGSMITH_TENANT_ID"); v != "" {
		opts.WorkspaceID = v
	}
	if v := os.Getenv("LANGSMITH_WORKSPACE_ID"); v != "" {
		opts.WorkspaceID = v
	}
	if flags.WorkspaceID != "" {
		opts.WorkspaceID = flags.WorkspaceID
	}
	switch {
	case flags.APIKey != "":
		opts.APIKey = flags.APIKey
	case os.Getenv("LANGSMITH_API_KEY") != "":
		if flags.Profile != "" {
			fmt.Fprintln(os.Stderr, "warning: --profile was specified, but LANGSMITH_API_KEY is set and takes precedence over saved profile auth")
		}
		opts.APIKey = os.Getenv("LANGSMITH_API_KEY")
	case hasProfile && (profile.AccessToken() != "" || profile.OAuth.RefreshToken != ""):
		// The SDK refreshes the profile's OAuth token on each request, under a
		// lock shared by every process using this config.
		opts.ProfileName = profileName
	case hasProfile && profile.HasAPIKey():
		key, err := profile.ResolveAPIKey()
		if err != nil {
			return opts, fmt.Errorf("profile %q: %w", profileName, err)
		}
		opts.APIKey = key
		// Route the resolved profile through the SDK (WithProfile) so an explicit
		// selection replaces the config's current_profile instead of inheriting
		// its tenant/base URL. APIKey is kept for the raw-HTTP helpers.
		opts.ProfileName = profileName
	}
	if cfgErr != nil && opts.APIKey == "" {
		return opts, cfgErr
	}

	return opts, nil
}

// ExitError prints an error to stderr and exits with exitcode.Input.
func ExitError(msg string) {
	exitWithCode(msg, exitcode.Input)
}

// ExitErrorf prints a formatted error to stderr and exits. The exit code comes
// from the first error among args (see ExitCode), or is exitcode.Input when
// there is none.
func ExitErrorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	for _, arg := range args {
		if err, ok := arg.(error); ok {
			fmt.Fprintln(os.Stderr, msg)
			exitProcess(ExitCode(err))
			return
		}
	}
	exitWithCode(msg, exitcode.Input)
}

// notFoundf returns an error for a LangSmith resource that a lookup did not
// find, which exits with exitcode.NotFound.
func notFoundf(format string, args ...any) error {
	return exitcode.WithCode(fmt.Errorf(format, args...), exitcode.NotFound)
}

// exitNotFoundf prints a not-found message and exits with the code for
// exitcode.NotFound.
func exitNotFoundf(format string, args ...any) {
	err := notFoundf(format, args...)
	fmt.Fprintln(os.Stderr, err)
	exitProcess(ExitCode(err))
}

// exitProcess ends the process. Tests replace it.
var exitProcess = os.Exit

func exitWithCode(msg string, code int) {
	fmt.Fprintln(os.Stderr, msg)
	exitProcess(code)
}
