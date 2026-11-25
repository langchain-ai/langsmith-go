package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ReportedError is an error a command has already reported in its own output
// format. The CLI exits with a failure status without printing it again.
type ReportedError struct {
	Err error
}

func (e ReportedError) Error() string { return e.Err.Error() }

func (e ReportedError) Unwrap() error { return e.Err }

// errorsToStderr prints a command's error to stderr and returns it as a
// ReportedError, so the error stays off stdout, where agents read JSON output.
func errorsToStderr(run func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := run(cmd, args)
		if err == nil {
			return nil
		}
		fmt.Fprintln(cmd.ErrOrStderr(), err)
		return ReportedError{Err: err}
	}
}
