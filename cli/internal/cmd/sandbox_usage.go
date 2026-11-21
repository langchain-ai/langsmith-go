package cmd

import (
	"context"
	"fmt"

	"github.com/langchain-ai/langsmith-go/cli/internal/cmdutil"
	"github.com/langchain-ai/langsmith-go/cli/internal/structured"
	"github.com/spf13/cobra"
)

var sandboxUsageCommand = structured.Command[struct{}]{
	Use:   "usage",
	Short: "Show the workspace's sandbox usage against its limits",
	Long: `Show the sandboxes, CPU and memory the workspace is using against its limits.
Only running sandboxes (provisioning or ready) count toward the limits; stopped
ones are counted separately and hold no CPU or memory. CPU is in millicores (1000m = 1 vCPU).

Examples:
  langsmith sandboxes usage
  langsmith sandboxes usage --format json --jq '.sandboxes_limit - .sandboxes_used'`,
	Args: cobra.NoArgs,
	Action: func(ctx context.Context, cmd *cobra.Command, in struct{}, args []string) (any, error) {
		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}
		resp, err := c.SDK.Sandboxes.GetUsage(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting sandbox usage: %w", err)
		}
		return resp, nil
	},
	Render: structured.PropertyList{
		Properties: []structured.Property{
			{Label: "Running", Template: "{{.SandboxesUsed}} / {{.SandboxesLimit}}"},
			{Label: "Stopped", Template: "{{.SandboxStoppedCount}}"},
			{Label: "CPU", Template: "{{.CPUUsedMillicores}}m / {{.CPULimitMillicores}}m"},
			{Label: "Memory", Template: "{{formatBytes .MemoryUsedBytes}} / {{formatBytes .MemoryLimitBytes}}"},
			{Label: "Snapshots", Template: "{{.SnapshotsCount}}"},
			{Label: "Registries", Template: "{{.RegistriesCount}}"},
		},
	},
}
