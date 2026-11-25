// Package cli is the langsmith command-line tool. Its commands live in this
// module next to the SDK they call, so one change can add an SDK method and
// the command that uses it.
//
// Build the binary from github.com/langchain-ai/langsmith-go/cmd/langsmith.
// Released binaries come from github.com/langchain-ai/langsmith-cli, which
// calls [Main] from its own main package.
package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/langchain-ai/langsmith-go/cli/internal/cmd"
)

// Main runs the langsmith CLI with the process arguments and returns its exit
// code. version, commit and date describe the build for `langsmith --version`.
func Main(version, commit, date string) int {
	rootCmd := cmd.NewRootCmd(version, fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date))
	if err := rootCmd.Execute(); err != nil {
		var reported cmd.ReportedError
		if !errors.As(err, &reported) {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		return cmd.ExitCode(err)
	}
	return 0
}
