package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/langchain-ai/langsmith-go/cli/internal/exitcode"
)

// ExitCode returns the process exit code for a command's error and, when the
// detailed codes are not enabled, writes a deprecation warning to stderr for a
// failure whose code changes in exitcode.DefaultIn.
func ExitCode(err error) int {
	return processExitCode(os.Stderr, err, GetFormat() == "json" || generatedJSON)
}

// generatedJSON reports whether the generated command that ran asked for JSON
// output with its own unparsed flags.
var generatedJSON bool

// processExitCode returns exitcode.For(err) when the detailed codes are enabled
// (see exitcode.EnvVar). Otherwise every failure exits exitcode.Input, and a
// failure whose detailed code differs gets a warning on w naming that code: a
// JSON object when asJSON, otherwise a "warning:" line.
func processExitCode(w io.Writer, err error, asJSON bool) int {
	code := exitcode.For(err)
	if exitcode.Detailed() || code == exitcode.OK || code == exitcode.Input {
		return code
	}
	warning := deprecationWarning{
		Deprecated:  true,
		Kind:        "exit_code",
		Name:        "1",
		RemovedIn:   exitcode.DefaultIn,
		Replacement: strconv.Itoa(code),
		Message: fmt.Sprintf("this failure exits 1, and from %s exits %d (%s); set %s=detailed to use the new exit codes now.",
			exitcode.DefaultIn, code, exitcode.Name(code), exitcode.EnvVar),
	}
	if asJSON {
		b, _ := json.Marshal(warning)
		fmt.Fprintln(w, string(b))
	} else {
		fmt.Fprintln(w, "warning: "+warning.Message)
	}
	return exitcode.Input
}
