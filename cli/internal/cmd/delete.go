package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// confirmsAnnotation marks a hand-written command that deletes or overwrites
// data and asks first through confirmDelete. `docs agents` marks these
// commands *(confirms)*, and TestDestructiveCommandsRequireYes runs each one.
const confirmsAnnotation = "langsmith/confirms"

// addYesFlag registers --yes, which skips confirmDelete, and marks cmd as a
// command that confirms.
func addYesFlag(cmd *cobra.Command, yes *bool, usage string) {
	cmd.Flags().BoolVar(yes, "yes", false, usage)
	markConfirms(cmd)
}

// markConfirms marks cmd as a command that confirms through confirmDelete.
// Use addYesFlag unless the skip flag is not --yes.
func markConfirms(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[confirmsAnnotation] = "true"
}

// yesInput is the Input of a structured command that confirms: it registers
// --yes and returns its value.
func yesInput(cmd *cobra.Command) *bool {
	var yes bool
	addYesFlag(cmd, &yes, "Skip confirmation prompt")
	return &yes
}

// requireYesWithoutTerminal makes confirmDelete refuse whenever stdin is not
// a terminal. Scripts rely on some of these commands running with no prompt
// and on a piped "y" confirming the rest, and the CLI's breaking-change policy
// keeps both working, with a deprecation warning, until the major version in
// strictConfirmationVersion. Set it to true in that release.
const requireYesWithoutTerminal = false

const strictConfirmationVersion = "1.0.0"

type deleteConfirmation struct {
	// target is what the command removes, as in "This permanently deletes
	// <target>".
	target   string
	identity string
	// verb replaces "deletes" in the warning, for commands that overwrite.
	verb string
	// skipFlag is the flag that skips the prompt; "" means --yes.
	skipFlag string
	// noSkip is set for commands that offer no way to skip the prompt.
	noSkip bool
	// unprompted is set where the command used to run with no prompt at all.
	// Until requireYesWithoutTerminal, it still runs without a terminal and
	// prints a deprecation warning.
	unprompted bool
	// acceptsUnterminated is set where the command's earlier prompt took a
	// "y" cut off by end of input. Until requireYesWithoutTerminal, a piped
	// answer still confirms these without a trailing newline; every other
	// command needs a complete line.
	acceptsUnterminated bool
	// jsonOutput selects a JSON deprecation warning for commands whose
	// --format and --jq are not parsed into the global flags.
	jsonOutput bool
}

// confirmDelete asks before a command deletes or overwrites data. When stdin
// is not a terminal it does not prompt: an agent or a script must pass the
// skip flag after the user approves. Until requireYesWithoutTerminal, a piped
// "y" (or, for an unprompted command, nothing) still confirms, with a
// deprecation warning.
func confirmDelete(cmd *cobra.Command, confirmation deleteConfirmation) error {
	verb := confirmation.verb
	if verb == "" {
		verb = "deletes"
	}
	skipFlag := confirmation.skipFlag
	if skipFlag == "" {
		skipFlag = "--yes"
	}
	stderr := cmd.ErrOrStderr()
	fmt.Fprintf(stderr, "WARNING: This permanently %s %s. This cannot be undone.\n", verb, confirmation.target)
	if confirmation.identity != "" {
		fmt.Fprintln(stderr, confirmation.identity)
	}

	if !inputIsTerminal(cmd.InOrStdin()) {
		if !requireYesWithoutTerminal && (confirmation.unprompted || readYes(cmd, confirmation.acceptsUnterminated)) {
			warnConfirmedWithoutTerminal(cmd, confirmation, skipFlag)
			return nil
		}
		if confirmation.noSkip {
			return errors.New("refusing to continue without confirmation: stdin is not a terminal and this command has no flag to skip the prompt. AI agents: stop and ask the user to run it in a terminal")
		}
		return fmt.Errorf("refusing to continue without confirmation: stdin is not a terminal. AI agents: stop and ask the user; only after they approve this specific change, rerun the command with %s", skipFlag)
	}

	fmt.Fprintln(stderr, "AI agents: do not answer this prompt. Stop and raise it to the user.")
	fmt.Fprint(stderr, "Continue? [y/N] ")
	if !readYes(cmd, false) {
		return errors.New("aborted: not confirmed")
	}
	return nil
}

// readYes reads one line from stdin and reports whether it is "y" or "yes".
// An answer cut off by end of input counts only when allowEOF is set; see
// deleteConfirmation.acceptsUnterminated.
func readYes(cmd *cobra.Command, allowEOF bool) bool {
	answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil && !allowEOF {
		return false
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes"
}

// warnConfirmedWithoutTerminal writes the deprecation warning for a change
// confirmed while stdin is not a terminal, as JSON under --format json or
// --jq like the other deprecation warnings.
func warnConfirmedWithoutTerminal(cmd *cobra.Command, confirmation deleteConfirmation, skipFlag string) {
	w := deprecationWarning{
		Deprecated:  true,
		Kind:        "behavior",
		Name:        "confirmation without a terminal",
		RemovedIn:   strictConfirmationVersion,
		Replacement: skipFlag,
		Message: fmt.Sprintf("running this command without a terminal and without %s is deprecated; from %s it will refuse. Pass %s now.",
			skipFlag, strictConfirmationVersion, skipFlag),
	}
	if confirmation.noSkip {
		w.Replacement = ""
		w.Message = fmt.Sprintf("confirming this command without a terminal is deprecated; from %s it will refuse. Run it in a terminal.",
			strictConfirmationVersion)
	}
	if confirmation.jsonOutput || GetFormat() == "json" {
		b, _ := json.Marshal(w)
		fmt.Fprintln(cmd.ErrOrStderr(), string(b))
		return
	}
	fmt.Fprintln(cmd.ErrOrStderr(), "warning: "+w.Message)
}
