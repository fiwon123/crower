package system

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/shell"
)

// Creator creates files and folders with the platform's own tools.
type Creator struct {
	shell *shell.Runner
}

// NewCreator builds a Creator that runs commands on the host.
func NewCreator() *Creator {
	return &Creator{shell: shell.New()}
}

// File creates one file.
func (c *Creator) File(path string) error {
	commandLine, err := CreateFileCommand(runtime.GOOS, path)
	if err != nil {
		return err
	}

	return c.run(commandLine, "create file", path)
}

// Folder creates one folder, including any missing parent.
func (c *Creator) Folder(path string) error {
	commandLine, err := CreateFolderCommand(runtime.GOOS, path)
	if err != nil {
		return err
	}

	return c.run(commandLine, "create folder", path)
}

// run executes one platform command and wraps its output into the error.
func (c *Creator) run(commandLine, operation, path string) error {
	out, err := c.shell.Exec(commandLine)
	if err != nil {
		return fmt.Errorf("%s %s: %w (output: %s)", operation, path, err, strings.TrimSpace(out))
	}

	return nil
}

// CreateFileCommand builds the platform's empty-file invocation.
//
// It is a pure function of goos so a table test can cover every platform the
// release ships (#129).
func CreateFileCommand(goos, path string) (string, error) {
	switch goos {
	case "windows":
		return "type nul > " + shell.QuoteFor(goos, path), nil
	case "linux", "darwin":
		return "touch " + shell.QuoteFor(goos, path), nil
	default:
		return "", fmt.Errorf("%w: create file on %q", errors.ErrUnsupportedOS, goos)
	}
}

// CreateFolderCommand builds the platform's folder invocation.
func CreateFolderCommand(goos, path string) (string, error) {
	switch goos {
	case "windows":
		return "mkdir " + shell.QuoteFor(goos, path), nil
	case "linux", "darwin":
		return "mkdir -p " + shell.QuoteFor(goos, path), nil
	default:
		return "", fmt.Errorf("%w: create folder on %q", errors.ErrUnsupportedOS, goos)
	}
}
