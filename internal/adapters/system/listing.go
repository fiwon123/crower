package system

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/fiwon123/crower/internal/application/listcontent"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/shell"
	"github.com/fiwon123/crower/pkg/utils"
)

// Listings runs the platform listings and reads a directory for a structured
// result.
type Listings struct {
	shell *shell.Runner

	// goos is the platform the shell command is built for. It is a field rather
	// than a read of runtime.GOOS inside each method so the dispatch can be
	// exercised by a test: runtime.GOOS cannot be overridden at run time.
	goos string

	// listProcesses is a field so a test can answer it without a process table.
	listProcesses func(partName string, skipUnknown bool) ([]utils.ProcessInfo, error)
}

// NewListings builds a Listings that reads the host.
func NewListings() *Listings {
	return &Listings{
		shell:         shell.New(),
		goos:          runtime.GOOS,
		listProcesses: utils.ListProcesses,
	}
}

// FolderText runs the platform's folder listing.
func (l *Listings) FolderText(path string) (string, error) {
	commandLine, err := FolderListingCommand(l.goos, path)
	if err != nil {
		return "", err
	}

	return l.shell.Exec(commandLine)
}

// Variables returns the raw environment listing.
func (l *Listings) Variables() (string, error) {
	commandLine, err := VariablesListingCommand(l.goos)
	if err != nil {
		return "", err
	}

	return l.shell.Exec(commandLine)
}

// Path returns the raw PATH value.
func (l *Listings) Path() (string, error) {
	commandLine, err := PathListingCommand(l.goos)
	if err != nil {
		return "", err
	}

	return l.shell.Exec(commandLine)
}

// FolderEntries reads a folder for a structured result.
//
// It is not a drop-in replacement for the text form: the path is used literally
// (no shell expansion), a symlink is followed, and a path that is not a
// directory is an error.
func (l *Listings) FolderEntries(path string) ([]listcontent.FolderEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("list folder %s: %w", path, err)
	}

	rows := make([]listcontent.FolderEntry, 0, len(entries))

	for _, entry := range entries {
		row := listcontent.FolderEntry{Name: entry.Name(), Dir: entry.IsDir()}

		info, err := entry.Info()
		if err != nil {
			// The entry vanished between the read and this call, or the OS
			// refused its metadata; the name and type are still real.
			row.Error = err.Error()
		} else {
			row.Size = info.Size()
			row.Mode = info.Mode().String()
			row.Modified = info.ModTime().Format(time.RFC3339)
		}

		rows = append(rows, row)
	}

	return rows, nil
}

// Processes lists the running processes, mapped to the result's row type.
func (l *Listings) Processes(partName string) ([]listcontent.ProcessRow, error) {
	processes, err := l.listProcesses(partName, true)
	if err != nil {
		return nil, err
	}

	rows := make([]listcontent.ProcessRow, 0, len(processes))

	for _, process := range processes {
		rows = append(rows, listcontent.ProcessRow{
			PID:  process.PID,
			Name: process.Name,
			Path: process.Path,
		})
	}

	return rows, nil
}

// The listing commands are ordinary shell invocations that behave the same on
// every POSIX system, so linux and darwin share one arm. Each function has a
// default arm that returns an error: an unsupported platform must not look like
// a successful empty listing.
func FolderListingCommand(goos, folderPath string) (string, error) {
	switch goos {
	case "windows":
		return "dir " + shell.QuoteFor(goos, folderPath), nil
	case "linux", "darwin":
		// -- stops a folder named like an option from changing the listing.
		return "ls -la -- " + shell.QuoteFor(goos, folderPath), nil
	default:
		return "", fmt.Errorf("%w: %q", errors.ErrUnsupportedOS, goos)
	}
}

func VariablesListingCommand(goos string) (string, error) {
	switch goos {
	case "windows":
		return "set", nil
	case "linux", "darwin":
		return "printenv", nil
	default:
		return "", fmt.Errorf("%w: %q", errors.ErrUnsupportedOS, goos)
	}
}

func PathListingCommand(goos string) (string, error) {
	switch goos {
	case "windows":
		// Quotes keep an entry containing "&" from becoming a separator.
		return `echo "%PATH%"`, nil
	case "linux", "darwin":
		return `echo "$PATH"`, nil
	default:
		return "", fmt.Errorf("%w: %q", errors.ErrUnsupportedOS, goos)
	}
}
