//go:build !linux && !darwin

package shell

import (
	"os/exec"
	"time"
)

// Isolate makes cancellation kill the process itself.
//
// There is no portable process group here, so a grandchild can outlive the
// shell. WaitDelay bounds the wait for the output pipe that grandchild may
// still hold; without it, a timed step could hang past its deadline.
func Isolate(command *exec.Cmd) {
	command.Cancel = func() error {
		if command.Process == nil {
			return nil
		}

		return command.Process.Kill()
	}

	command.WaitDelay = time.Second
}
