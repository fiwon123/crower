//go:build windows

package cmd

import (
	"os"
	osexec "os/exec"
)

// startReplacement spawns the updated binary. Windows has no exec-in-place, so
// the new process is started independently of the closing window.
func startReplacement(executable string, args []string) error {
	executable = cleanExecutable(executable)

	command := osexec.Command(executable, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Start()
}
