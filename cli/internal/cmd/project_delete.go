package cmd

import (
	"fmt"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/spf13/cobra"
)

func newProjectDeleteCmd() *cobra.Command {
	var projectName, projectID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Permanently delete a tracing project and all of its traces",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			id, err := resolveSessionID(ctx, c, projectName, projectID, "project delete")
			if err != nil {
				return err
			}
			project, err := c.SDK.Sessions.GetV2(ctx, id, langsmith.SessionGetV2Params{
				IncludeStats: langsmith.F(true),
			})
			if err != nil {
				return fmt.Errorf("getting tracing project %s: %w", id, err)
			}

			if err := confirmDelete(cmd, deleteConfirmation{
				target:   "the tracing project and all of its traces",
				identity: fmt.Sprintf("Project: %q (id: %s, runs: %d)", project.Name, project.ID, project.RunCount),
				noSkip:   true,
			}); err != nil {
				return err
			}

			if err := c.SDK.Sessions.Delete(ctx, id); err != nil {
				return fmt.Errorf("deleting tracing project %s: %w", id, err)
			}
			// The API accepts the delete (202) and removes the project and
			// its traces in the background, so it is not gone yet.
			return output.OutputJSON(map[string]any{
				"status": "deleting",
				"id":     project.ID,
				"name":   project.Name,
			}, "")
		},
	}

	addProjectFlags(cmd, &projectName, &projectID)
	// Deleting a project and its traces has no --yes: only a person at a
	// terminal can confirm it.
	markConfirms(cmd)
	return cmd
}
