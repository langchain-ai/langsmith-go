package cmd

import "github.com/spf13/cobra"

type sandboxMessage struct {
	Name    string `json:"name,omitempty"`
	Message string `json:"message"`
}

func newSandboxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sandboxes",
		Aliases: []string{"sandbox"},
		Short:   "Manage and interact with sandboxes (experimental)",
		Long: `Manage and interact with sandboxes (currently in experimental preview).

Common workflows:

  # Build a snapshot from a Docker image
  langsmith sandboxes snapshot build my-snap --docker-image ubuntu:24.04

  # Create a sandbox from the snapshot and wait until it is ready
  langsmith sandboxes create my-vm --snapshot-id <id> --wait

  # Run a command inside the sandbox
  langsmith sandboxes exec my-vm -- uname -a

  # Copy files in and out
  langsmith sandboxes cp ./app.tar.gz my-vm:/tmp/
  langsmith sandboxes cp my-vm:/tmp/report.pdf .

  # Keep memory across stops, so a restart resumes running processes
  langsmith sandboxes create my-vm --preserve-memory-on-stop

  # Capture a snapshot with memory, so new sandboxes resume from it
  langsmith sandboxes snapshot capture warm --from my-vm --include-memory --wait

  # Open an interactive shell
  langsmith sandboxes console my-vm

  # Tunnel a remote port (e.g. Postgres) to localhost
  langsmith sandboxes tunnel my-vm --remote-port 5432

  # Generate an authenticated HTTP URL for a service inside the sandbox
  langsmith sandboxes service-url my-vm --port 8000

  # Generate a credential-free link that downloads one file from the sandbox
  langsmith sandboxes generate-download-url my-vm --path /tmp/report.pdf

  # Set up SSH access (writes ~/.ssh/config so "ssh sandbox-my-vm" works)
  # Requires sshd to be installed in your Docker image.
  langsmith sandboxes ssh-setup my-vm`,
	}

	// Lifecycle
	cmd.AddCommand(sandboxCreateCommand.Cobra())
	cmd.AddCommand(sandboxListCommand.Cobra())
	cmd.AddCommand(sandboxGetCommand.Cobra())
	cmd.AddCommand(sandboxUpdateCommand.Cobra())
	cmd.AddCommand(sandboxDeleteCommand.Cobra())
	cmd.AddCommand(sandboxUsageCommand.Cobra())
	cmd.AddCommand(sandboxStartCommand.Cobra())
	cmd.AddCommand(sandboxStopCommand.Cobra())

	// Connectivity
	cmd.AddCommand(newSandboxExecCmd())
	cmd.AddCommand(newSandboxConsoleCmd())
	cmd.AddCommand(sandboxServiceURLCommand.Cobra())
	cmd.AddCommand(sandboxDownloadURLCommand.Cobra())
	cmd.AddCommand(sandboxCopyCommand.Cobra())
	cmd.AddCommand(sandboxListFilesCommand.Cobra())
	cmd.AddCommand(newSandboxTunnelCmd())
	cmd.AddCommand(newSandboxSSHSetupCmd())

	// Sub-resources
	cmd.AddCommand(sandboxSnapshotCommand.Cobra())

	return cmd
}
