package cmd

import (
	"context"
	"fmt"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/spf13/cobra"
)

func newHubDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete [OWNER/]REPO",
		Short: "Delete an agent or skill repo (and its owned child repos)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, name, _, err := parseHubOwnerRepo(args[0])
			if err != nil {
				return err
			}
			if !yes {
				if err := confirmDelete(cmd, deleteConfirmation{
					acceptsUnterminated: true,
					target:              "the hub repo and its owned child repos",
					identity:            fmt.Sprintf("Repo: %s/%s", owner, name),
				}); err != nil {
					return err
				}
			}
			c := MustGetClient()
			ctx := context.Background()

			// No repo_type filter: `hub delete` deletes the repo whatever its type.
			if err := c.SDK.Repos.Directories.Delete(ctx, owner, name, langsmith.RepoDirectoryDeleteParams{}); err != nil {
				return fmt.Errorf("deleting %s/%s: %w", owner, name, err)
			}
			if err := output.OutputJSON(map[string]any{
				"status": "deleted",
				"owner":  owner,
				"repo":   name,
			}, ""); err != nil {
				ExitErrorf("%v", err)
			}
			return nil
		},
	}
	addYesFlag(cmd, &yes, "Skip confirmation prompt")
	return cmd
}
