//go:build linux || darwin

package shell

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Isolate puts the child in its own process group and makes cancellation kill
// the group.
//
// Killing the shell alone is not enough: `sh -c "sleep 5"` leaves sleep holding
// the output pipe, so CombinedOutput waits for it even after the shell is dead.
// A timeout must end the step, not the shell that started it.
//
// Call this only for a command that can be cancelled: a new process group no
// longer receives the terminal's Ctrl+C, so an ordinary command keeps the
// process group it has always had. A timed step is still bounded by its
// timeout while crower runs; if crower itself is killed, the group can outlive
// it.
func Isolate(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	command.Cancel = func() error { return killProcessGroup(command) }

	// A survivor that left the group (setsid) can still hold the output pipe,
	// and Wait would block on it forever. The delay bounds that; the deadline
	// has already expired, so the caller still reports a timeout.
	command.WaitDelay = time.Second
}

func killProcessGroup(command *exec.Cmd) error {
	if command.Process == nil {
		return nil
	}

	if err := syscall.Kill(-command.Process.Pid, syscall.SIGKILL); err != nil {
		// The process is already gone: the ordinary "done" answer, not a
		// cancellation failure exec would report instead of the timeout.
		if errors.Is(err, syscall.ESRCH) {
			return os.ErrProcessDone
		}

		return err
	}

	return nil
}
