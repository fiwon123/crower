package system

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/shell"
)

// Remover deletes files and folders with the platform's own tools.
type Remover struct {
	shell *shell.Runner
}

// NewRemover builds a Remover that runs commands on the host.
func NewRemover() *Remover {
	return &Remover{shell: shell.New()}
}

// File deletes one file.
func (r *Remover) File(path string) error {
	return r.remove(path, false)
}

// Folder deletes one folder and everything in it.
func (r *Remover) Folder(path string) error {
	return r.remove(path, true)
}

// remove runs the platform's delete command.
func (r *Remover) remove(path string, folder bool) error {
	commandLine, err := RemoveCommand(runtime.GOOS, path, folder)
	if err != nil {
		return err
	}

	out, err := r.shell.Exec(commandLine)
	if err != nil {
		return fmt.Errorf("delete %s: %w (output: %s)", path, err, strings.TrimSpace(out))
	}

	return nil
}

// RemoveCommand builds the platform's delete invocation.
//
// It is a pure function of goos so a table test can cover every platform the
// release ships: runtime.GOOS cannot be overridden at run time, and an arm no
// test exercises is how macOS spent a release reporting success without
// copying anything (#129).
func RemoveCommand(goos, path string, folder bool) (string, error) {
	switch goos {
	case "windows":
		if folder {
			return "rmdir /s /q " + shell.QuoteFor(goos, path), nil
		}

		return "del " + shell.QuoteFor(goos, path), nil

	case "linux", "darwin":
		commandLine := "rm "
		if folder {
			commandLine = "rm -r "
		}

		return commandLine + shell.QuoteFor(goos, path), nil

	default:
		return "", fmt.Errorf("%w: delete on %q", errors.ErrUnsupportedOS, goos)
	}
}
