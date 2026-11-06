//go:build !windows

package cmd

import (
	"os"
	"syscall"
)

// startReplacement replaces this process image with the binary at the path
// captured before the update. exec keeps the PID and its service/session, so a
// systemd or launchd scope cannot kill the replacement when the old process
// exits; the window is already closed by the caller.
func startReplacement(executable string, args []string) error {
	executable = cleanExecutable(executable)

	return syscall.Exec(executable, replacementArgv(executable, args), os.Environ())
}
