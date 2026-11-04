// Package shell runs one command line through the platform's shell.
//
// It exists so a domain can run `ls`, `cp` or `xdg-open` without depending on
// the execute operation. A domain owns an operation, its history and its notes;
// this is infrastructure every operation may need, and it carries no state.
package shell

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fiwon123/crower/internal/errors"
)

// Runner starts the platform's shell for a command line.
//
// It has no state, so any number of domains can hold one.
type Runner struct{}

// New returns a Runner for the running platform.
func New() *Runner {
	return &Runner{}
}

// Exec runs line and returns its combined output.
//
// The whole line is handed to the shell as a single argument. It used to be
// split on whitespace and passed as separate argv entries after -c, which makes
// the shell treat everything after the first word as positional parameters:
// `echo hello` ran `echo` with no arguments, and `ls /nope` listed the working
// directory instead of failing.
func (r *Runner) Exec(line string) (string, error) {
	c, err := r.Command(line)
	if err != nil {
		return "", err
	}

	out, err := c.CombinedOutput()
	return string(out), err
}

// Start starts line without waiting for it.
//
// A browser or a file manager is expected to outlive the command that opened
// it, so waiting would hold the terminal until the window closed.
func (r *Runner) Start(line string) error {
	c, err := r.Command(line)
	if err != nil {
		return err
	}

	return c.Start()
}

// Interactive runs name with arg attached to the terminal itself.
//
// An editor needs the terminal, so its streams cannot be captured: it either
// has a tty or it is unusable.
func (r *Runner) Interactive(name string, arg string) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("interactive terminal is not supported on windows")
	}

	c := exec.Command(name, arg)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

// Command builds the platform's shell invocation for one command line, for a
// caller that needs the process itself (for example to set its environment).
func (r *Runner) Command(line string) (*exec.Cmd, error) {
	return commandFor(runtime.GOOS, line)
}

// CommandFor builds the invocation of the shell a command names.
//
// An empty name is the platform's own shell, which is what every command did
// before the portable schema. A named shell runs on every platform: "sh" is
// POSIX sh, "pwsh" is PowerShell, and "cmd" is the Windows command interpreter.
// The name is a portable declaration, not a platform guess -- a command that
// asks for pwsh runs pwsh on Linux too, and fails with the OS's own "not found"
// when it is not installed.
func (r *Runner) CommandFor(shellName, line string) (*exec.Cmd, error) {
	return commandForShell(runtime.GOOS, shellName, line)
}

// commandFor is the platform choice, with goos as a parameter so a test can
// exercise every arm: runtime.GOOS cannot be overridden at run time.
func commandFor(goos string, line string) (*exec.Cmd, error) {
	return commandForShell(goos, "", line)
}

// commandForShell is the platform default plus the three shells a command may
// name explicitly.
//
// pwsh runs with -NoProfile: a shared command should behave the same on every
// machine, and a user's profile can change what a command line means.
func commandForShell(goos, shellName, line string) (*exec.Cmd, error) {
	switch shellName {
	case "":
		switch goos {
		case "windows":
			return exec.Command("cmd", "/c", line), nil
		case "linux", "darwin":
			return exec.Command("sh", "-c", line), nil
		default:
			return nil, fmt.Errorf("%w: %q", errors.ErrUnsupportedOS, goos)
		}
	case "sh":
		return exec.Command("sh", "-c", line), nil
	case "pwsh":
		return exec.Command("pwsh", "-NoProfile", "-Command", line), nil
	case "cmd":
		return exec.Command("cmd", "/c", line), nil
	default:
		return nil, fmt.Errorf("%w: %q", errors.ErrUnknownShell, shellName)
	}
}
