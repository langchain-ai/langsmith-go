//go:build windows

package cmd

// terminalDeliversInterrupt is always true on Windows: a console Ctrl-C
// reaches every process attached to the console, and os.Process cannot send
// an interrupt to another process.
func terminalDeliversInterrupt() bool { return true }
