package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/spf13/cobra"
)

// nullableStringInput distinguishes an omitted JSON field from an explicit
// null, which is needed to clear a schedule with PATCH.
type nullableStringInput struct {
	Set   bool
	Value *string
}

func (v *nullableStringInput) UnmarshalJSON(data []byte) error {
	v.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		v.Value = nil
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	v.Value = &value
	return nil
}

func applyInsightConfigMetadataToNewParams(input insightConfigFile, params *langsmith.SessionInsightConfigNewParams) {
	if input.Description != nil {
		params.Description = langsmith.F(*input.Description)
	}
	if input.ScheduleCron.Set {
		if input.ScheduleCron.Value == nil {
			params.ScheduleCron = langsmith.Null[string]()
		} else {
			params.ScheduleCron = langsmith.F(strings.TrimSpace(*input.ScheduleCron.Value))
		}
	}
}

func applyInsightConfigMetadataToUpdateParams(input insightConfigFile, params *langsmith.SessionInsightConfigUpdateParams) {
	if input.Description != nil {
		params.Description = langsmith.F(*input.Description)
	}
	if input.ScheduleCron.Set {
		if input.ScheduleCron.Value == nil {
			params.ScheduleCron = langsmith.Null[string]()
		} else {
			params.ScheduleCron = langsmith.F(strings.TrimSpace(*input.ScheduleCron.Value))
		}
	}
}

func newInsightsConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Insights job configurations",
		Long: `Manage the saved configurations used to run Insights jobs.

Configurations contain the analysis parameters and may include a recurring
schedule. Updating a configuration does not start a job; use 'insights run'
to run an existing configuration immediately.`,
	}
	cmd.AddCommand(newInsightsConfigListCmd())
	cmd.AddCommand(newInsightsConfigUpdateCmd())
	cmd.AddCommand(newInsightsConfigDeleteCmd())
	return cmd
}

func newInsightsConfigListCmd() *cobra.Command {
	var (
		project          string
		projectID        string
		includePrebuilts bool
		outputFile       string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Insights job configurations for a project",
		Example: `  langsmith insights config list --project my-app
  langsmith insights config list --project-id PROJECT_ID --include-prebuilt --format json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			sessionID, err := resolveSessionID(cmd.Context(), c, project, projectID, "insights config list")
			if err != nil {
				return err
			}
			configs, err := c.SDK.Sessions.Insights.Configs.List(cmd.Context(), sessionID, langsmith.SessionInsightConfigListParams{
				IncludePrebuilts: langsmith.F(includePrebuilts),
			})
			if err != nil {
				return fmt.Errorf("listing insight configs: %w", err)
			}

			if GetFormat() == "pretty" {
				rows := make([][]string, 0, len(*configs))
				for _, config := range *configs {
					schedule := "none"
					if !config.JSON.ScheduleCron.IsNull() && config.ScheduleCron != "" {
						schedule = config.ScheduleCron
					}
					rows = append(rows, []string{
						config.Name,
						config.ID,
						schedule,
						fmt.Sprintf("%t", config.Prebuilt),
					})
				}
				output.OutputTable([]string{"Name", "ID", "Schedule", "Prebuilt"}, rows, "Insights Job Configurations")
				return nil
			}

			data := make([]json.RawMessage, 0, len(*configs))
			for _, config := range *configs {
				data = append(data, json.RawMessage(config.JSON.RawJSON()))
			}
			return output.OutputJSON(data, outputFile)
		},
	}

	addProjectFlags(cmd, &project, &projectID)
	cmd.Flags().BoolVar(&includePrebuilts, "include-prebuilt", false, "Include built-in configurations")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write JSON output to a file")
	return cmd
}

func newInsightsConfigUpdateCmd() *cobra.Command {
	var (
		project    string
		projectID  string
		configFile string
		outputFile string
	)

	cmd := &cobra.Command{
		Use:   "update CONFIG_ID",
		Short: "Update an Insights job configuration",
		Long: `Update an Insights job configuration from a manual-mode JSON file.

The file replaces the saved name and analysis configuration. Omitting description
preserves its current value; set it to an empty string to clear its contents.
Omitting schedule_cron preserves its current value; set it to null to disable
the schedule.

When selecting models explicitly, cluster_model and summary_model must be set
together. Each value must be openai, anthropic, or the UUID of a workspace
model configuration; model display names are not accepted.

Updating a configuration does not start a job.`,
		Example: `  langsmith insights config update CONFIG_ID --project my-app --config insights.json
  cat insights.json | langsmith insights config update CONFIG_ID --project my-app --config -`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := loadInsightConfigFile(configFile)
			if err != nil {
				return err
			}
			config, err := input.toSDKParams()
			if err != nil {
				return err
			}
			params := langsmith.SessionInsightConfigUpdateParams{
				Name:   langsmith.F(strings.TrimSpace(input.Name)),
				Config: langsmith.F(config),
			}
			applyInsightConfigMetadataToUpdateParams(input, &params)

			c, err := getClient()
			if err != nil {
				return err
			}
			sessionID, err := resolveSessionID(cmd.Context(), c, project, projectID, "insights config update")
			if err != nil {
				return err
			}
			updated, err := c.SDK.Sessions.Insights.Configs.Update(cmd.Context(), sessionID, args[0], params)
			if err != nil {
				return fmt.Errorf("updating insight config %s: %w", args[0], err)
			}

			if GetFormat() == "pretty" {
				fmt.Println("Insights job configuration updated")
				fmt.Println()
				fmt.Printf("Name:      %s\n", updated.Name)
				fmt.Printf("ID:        %s\n", updated.ID)
				schedule := "none"
				if !updated.JSON.ScheduleCron.IsNull() && updated.ScheduleCron != "" {
					schedule = updated.ScheduleCron
				}
				fmt.Printf("Schedule:  %s\n", schedule)
				return nil
			}
			return output.OutputJSON(json.RawMessage(updated.JSON.RawJSON()), outputFile)
		},
	}

	addProjectFlags(cmd, &project, &projectID)
	cmd.Flags().StringVar(&configFile, "config", "", `Configuration JSON file (use "-" for stdin)`)
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write JSON output to a file")
	_ = cmd.MarkFlagRequired("config")
	return cmd
}

func newInsightsConfigDeleteCmd() *cobra.Command {
	var (
		project   string
		projectID string
		yes       bool
	)

	cmd := &cobra.Command{
		Use:   "delete CONFIG_ID",
		Short: "Delete an Insights job configuration",
		Example: `  langsmith insights config delete CONFIG_ID --project my-app
  langsmith insights config delete CONFIG_ID --project-id PROJECT_ID --yes`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configID := args[0]
			c, err := getClient()
			if err != nil {
				return err
			}
			sessionID, err := resolveSessionID(cmd.Context(), c, project, projectID, "insights config delete")
			if err != nil {
				return err
			}
			if !yes {
				if err := confirmDelete(cmd, deleteConfirmation{
					target:   "the Insights job configuration and all associated jobs and reports",
					identity: fmt.Sprintf("Config ID: %s", configID),
				}); err != nil {
					return err
				}
			}

			deleted, err := c.SDK.Sessions.Insights.Configs.Delete(cmd.Context(), sessionID, configID)
			if err != nil {
				return fmt.Errorf("deleting insight config %s: %w", configID, err)
			}
			return output.OutputJSON(map[string]any{
				"id":      deleted.ID,
				"message": deleted.Message,
			}, "")
		},
	}

	addProjectFlags(cmd, &project, &projectID)
	addYesFlag(cmd, &yes, "Skip confirmation prompt")
	return cmd
}

func newInsightsRunCmd() *cobra.Command {
	var (
		project    string
		projectID  string
		wait       bool
		outputFile string
	)

	cmd := &cobra.Command{
		Use:   "run CONFIG_ID",
		Short: "Run an existing Insights job configuration now",
		Long: `Start an Insights job from an existing saved configuration.

The saved configuration and its schedule are not changed.`,
		Example: `  langsmith insights run CONFIG_ID --project my-app
  langsmith insights run CONFIG_ID --project-id PROJECT_ID --wait --format json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			sessionID, err := resolveSessionID(cmd.Context(), c, project, projectID, "insights run")
			if err != nil {
				return err
			}
			created, err := c.SDK.Sessions.Insights.New(cmd.Context(), sessionID, langsmith.SessionInsightNewParams{
				CreateRunClusteringJobRequest: langsmith.CreateRunClusteringJobRequestParam{
					ConfigID: langsmith.F(args[0]),
				},
			})
			if err != nil {
				return fmt.Errorf("running insight config %s: %w", args[0], err)
			}
			var waitErr error
			if wait {
				if waitErr, err = waitForCreatedInsight(cmd.Context(), c.SDK, sessionID, created); err != nil {
					return err
				}
			}

			if GetFormat() == "pretty" {
				printInsightCreatePretty(created, sessionID, args[0], wait && waitErr == nil, "")
				return waitErr
			}
			if err := output.OutputJSON(insightCreateResponseToMap(created, sessionID, args[0]), insightOutputFile(outputFile, waitErr)); err != nil {
				return err
			}
			return waitErr
		},
	}

	addProjectFlags(cmd, &project, &projectID)
	cmd.Flags().BoolVar(&wait, "wait", false, "Wait for the report to succeed or fail")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write JSON output to a file")
	return cmd
}
