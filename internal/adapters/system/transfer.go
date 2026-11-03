// Package system builds and runs the platform commands for file operations.
//
// It is the technology side: which executable copies a folder on Windows, how a
// path is quoted, and how the command is run. The use cases describe *what* to
// copy; this package knows *how* the host does it.
package system

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/shell"
)

// Transfer copies and moves paths with the platform's own tools.
type Transfer struct {
	shell *shell.Runner
}

// NewTransfer builds a Transfer that runs commands on the host.
func NewTransfer() *Transfer {
	return &Transfer{shell: shell.New()}
}

// Copy copies one path into a destination folder.
func (t *Transfer) Copy(source, destination string, folder bool) error {
	commandLine, err := CopyCommand(runtime.GOOS, source, destination, folder)
	if err != nil {
		return err
	}

	return t.run(commandLine, "copy", source)
}

// Move moves one path into a destination folder.
func (t *Transfer) Move(source, destination string) error {
	commandLine, err := MoveCommand(runtime.GOOS, source, destination)
	if err != nil {
		return err
	}

	return t.run(commandLine, "move", source)
}

// run executes one platform command and wraps its output into the error.
func (t *Transfer) run(commandLine, operation, source string) error {
	out, err := t.shell.Exec(commandLine)
	if err != nil {
		return fmt.Errorf("%s %s: %w (output: %s)", operation, source, err, strings.TrimSpace(out))
	}

	return nil
}

// CopyCommand builds the platform's copy invocation.
//
// It is a pure function of goos so a table test can cover every platform the
// release ships: runtime.GOOS cannot be overridden at run time, and an arm no
// test exercises is how macOS spent a release reporting success without copying
// anything (#129).
func CopyCommand(goos, source, destFolder string, folder bool) (string, error) {
	switch goos {
	case "windows":
		if folder {
			return "xcopy " + shell.QuoteFor(goos, source) + " " + shell.QuoteFor(goos, destFolder) + " /E /I", nil
		}

		return "copy " + shell.QuoteFor(goos, source) + " " + shell.QuoteFor(goos, destFolder), nil

	case "linux", "darwin":
		commandLine := "cp "
		if folder {
			commandLine += "-r "
		}

		return commandLine + shell.QuoteFor(goos, source) + " " + shell.QuoteFor(goos, destFolder), nil

	default:
		return "", fmt.Errorf("%w: copy on %q", errors.ErrUnsupportedOS, goos)
	}
}

// MoveCommand builds the platform's move invocation.
func MoveCommand(goos, source, destFolder string) (string, error) {
	switch goos {
	case "windows":
		return "move " + shell.QuoteFor(goos, source) + " " + shell.QuoteFor(goos, destFolder), nil
	case "linux", "darwin":
		return "mv " + shell.QuoteFor(goos, source) + " " + shell.QuoteFor(goos, destFolder), nil
	default:
		return "", fmt.Errorf("%w: move on %q", errors.ErrUnsupportedOS, goos)
	}
}
