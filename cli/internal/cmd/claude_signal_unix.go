//go:build !windows

package cmd

import (
	"os"

	"golang.org/x/sys/unix"
)

// terminalDeliversInterrupt reports whether a Ctrl-C reaches Claude Code
// directly: this process has a controlling terminal whose foreground process
// group is this process's, which Claude Code shares. The controlling terminal
// is checked rather than stdin, since a terminal's Ctrl-C reaches its
// foreground group however stdin is redirected. A SIGINT received otherwise
// was sent to this process alone.
func terminalDeliversInterrupt() bool {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return false
	}
	defer tty.Close()
	pgrp, err := unix.IoctlGetInt(int(tty.Fd()), unix.TIOCGPGRP)
	return err == nil && pgrp == unix.Getpgrp()
}
