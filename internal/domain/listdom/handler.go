package listdom

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/errors"

	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/pkg/utils"
)

// ErrUnsupportedOS is returned when a listing command has no shell invocation
// for the running platform.
func unsupportedOSError(goos string) error {
	return fmt.Errorf("%w: %q", errors.ErrUnsupportedOS, goos)
}

type Handler struct {
	app     *app.Config
	execute *executedom.Handler

	// goos is the platform the shell command is built for. It is a field
	// rather than a read of runtime.GOOS inside each method so the dispatch
	// can be exercised by a test: runtime.GOOS cannot be overridden at run
	// time.
	goos string
}

func NewHandler(app *app.Config) *Handler {

	execute := executedom.NewHandler(app)

	return &Handler{
		app:     app,
		execute: execute,
		goos:    runtime.GOOS,
	}
}

// CommandsTable renders the table `list` prints, header included.
//
// It returns text rather than printing it because the same table is also a
// prompt: delete, execute, update and reset show it before asking which row to
// use, and there it is a diagnostic. A prompt that wrote to the result stream
// put a table on stdout during a failed run, where a caller piping the results
// read it as data.
func (h *Handler) CommandsTable() string {
	return h.commandsTable(h.app.OrderKeys, h.app.AllCommandsByName)
}

func (h *Handler) commandsTable(orderKeys []string, allCommands command.MapData) string {
	var builder strings.Builder

	builder.WriteString("------------------------------------------------\n")
	fmt.Fprintf(&builder, "%-3s %-12s %-16s %s\n", "Row", "Name", "Aliases", "Exec")
	builder.WriteString("------------------------------------------------\n")

	for i, key := range orderKeys {
		single := allCommands.Get(key)
		if single == nil {
			continue
		}

		fmt.Fprintf(&builder, "%-3d %-12s %-16s %s\n", i, single.Name, strings.Join(single.AllAlias, ","), single.Exec)
	}

	return builder.String()
}

// FolderResult is the result of `list --folder --json`.
type FolderResult struct {
	Folder  string        `json:"folder"`
	Entries []FolderEntry `json:"entries"`
}

// FolderResult reads folderPath for a JSON listing.
//
// The text rendering is `ls -la` (or `dir`), whose columns depend on platform
// and locale, and parsing that back into fields would be guesswork. A directory
// read gives the JSON result the same entries on every platform and needs no
// shell arm at all. It is not a drop-in replacement for the text form: the path
// is used literally (no shell expansion), a symlink is followed, and a path
// that is not a directory is an error.
func (h *Handler) FolderResult(folderPath string) (FolderResult, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return FolderResult{}, fmt.Errorf("list folder %s: %w", folderPath, err)
	}

	payload := FolderResult{Folder: folderPath, Entries: make([]FolderEntry, 0, len(entries))}

	for _, entry := range entries {
		row := FolderEntry{Name: entry.Name(), Dir: entry.IsDir()}

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

		payload.Entries = append(payload.Entries, row)
	}

	return payload, nil
}

// ListProcess returns the processes running on the user operational system
// (OS), filtered by a case-insensitive part of the name.
func (h *Handler) ListProcess(args []string) ([]utils.ProcessInfo, error) {
	partName := ""
	if len(args) > 0 {
		partName = args[0]
	}

	return utils.ListProcesses(partName, true)
}

// ListHistory returns the history table.
func (h *Handler) ListHistory() string {
	return h.app.History.GetList()
}

// ListFolder lists the files and folders under folderPath.
func (h *Handler) ListFolder(folderPath string) (string, error) {
	commandLine, err := listFolderCommand(h.goos, folderPath)
	if err != nil {
		return "", err
	}

	return h.execute.PerformExecute(commandLine)
}

// ListSystem lists every system variable.
func (h *Handler) ListSystem() (string, error) {
	commandLine, err := listSystemCommand(h.goos)
	if err != nil {
		return "", err
	}

	return h.execute.PerformExecute(commandLine)
}

// ListSysPath lists the system path variable.
func (h *Handler) ListSysPath() (string, error) {
	commandLine, err := listSysPathCommand(h.goos)
	if err != nil {
		return "", err
	}

	return h.execute.PerformExecute(commandLine)
}

// The listing commands are ordinary shell invocations that behave the same on
// every POSIX system, so linux and darwin share one arm. Each function has a
// default arm that returns an error: an unsupported platform must not look like
// a successful empty listing.
func listFolderCommand(goos string, folderPath string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf(`dir "%s"`, folderPath), nil
	case "linux", "darwin":
		return fmt.Sprintf(`ls -la "%s"`, folderPath), nil
	default:
		return "", unsupportedOSError(goos)
	}
}

func listSystemCommand(goos string) (string, error) {
	switch goos {
	case "windows":
		return "set", nil
	case "linux", "darwin":
		return "printenv", nil
	default:
		return "", unsupportedOSError(goos)
	}
}

func listSysPathCommand(goos string) (string, error) {
	switch goos {
	case "windows":
		return "echo %PATH%", nil
	case "linux", "darwin":
		return `echo "$PATH"`, nil
	default:
		return "", unsupportedOSError(goos)
	}
}
