package cmd

import (
	"fmt"

	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/spf13/cobra"
)

func newAppsDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete APP_ID_OR_NAME",
		Short: "Delete a custom app by ID or name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := MustGetClient()
			ctx := cmd.Context()

			app, err := resolveCustomApp(ctx, c, args[0])
			if err != nil {
				return err
			}
			id, name := app.ID, app.Name

			if !yes {
				if err := confirmDelete(cmd, deleteConfirmation{
					acceptsUnterminated: true,
					target:              "the custom app",
					identity:            fmt.Sprintf("App: %q (id: %s)", name, id),
				}); err != nil {
					return err
				}
			}

			if err := c.RawDelete(ctx, c.CustomAppPath(id), nil); err != nil {
				// Org apps can't be deleted here.
				if client.IsNotFound(err) && app.tier() == appScopeOrg {
					return fmt.Errorf("custom app %q (%s) is shared with the organization, so this workspace can't delete it. Claim a workspace copy with `langsmith apps claim %s` and delete that, or ask an organization admin", name, id, id)
				}
				return fmt.Errorf("deleting custom app %s: %w", id, err)
			}
			return output.OutputJSON(map[string]any{
				"status": "deleted",
				"app_id": id,
				"name":   name,
			}, "")
		},
	}

	addYesFlag(cmd, &yes, "Skip confirmation prompt")
	return cmd
}
