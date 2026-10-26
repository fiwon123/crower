package listdom

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/command"

	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/pkg/utils"
)

// ErrUnsupportedOS is returned when a listing command has no shell invocation
// for the running platform.
var ErrUnsupportedOS = errors.New("unsupported operating system")

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

// List all commands in order
func (h *Handler) ListCommands() {
	h.app.Logger.Info("------------------------------------------------")
	h.print(h.app.OrderKeys, h.app.AllCommandsByName)
}

func (h *Handler) print(orderKeys []string, allCommands command.MapData) {
	h.app.Logger.Info(fmt.Sprintf("%-3s %-12s %-16s %-8s \n", "Row", "Name", "Aliases", "Exec"))
	h.app.Logger.Info("------------------------------------------------")

	for i, key := range orderKeys {
		command := allCommands.Get(key)
		h.app.Logger.Info(fmt.Sprintf("%-3d %-12s %-16v %-8s \n", i, command.Name, strings.Join(command.AllAlias, ","), command.Exec))
	}
}

// List all ListProcess running on user operational system (OS).
func (h *Handler) ListProcess(args []string) error {

	partName := ""
	if len(args) > 0 {
		partName = args[0]
	}

	out, err := utils.GetAllProcess(partName, true)
	if err != nil {
		h.app.Logger.Error("Error getting processes:", err)
		return err
	}

	h.app.Logger.Info(out)

	return nil

}

// List all history
func (h *Handler) ListHistory() error {

	h.app.Logger.Info(h.app.History.GetList())

	return nil
}

// List all files and folder from a folderpath
func (h *Handler) ListFolder(folderPath string) (string, error) {
	commandLine, err := listFolderCommand(h.goos, folderPath)
	if err != nil {
		return "", err
	}

	return h.execute.PerformExecute(commandLine)
}

// List all system variable
func (h *Handler) ListSystem() (string, error) {
	commandLine, err := listSystemCommand(h.goos)
	if err != nil {
		return "", err
	}

	return h.execute.PerformExecute(commandLine)
}

// List system path variable
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
		return fmt.Sprintf("dir '%s'", folderPath), nil
	case "linux", "darwin":
		return fmt.Sprintf("\"ls '%s'\"", folderPath), nil
	default:
		return "", unsupportedOSError(goos)
	}
}

func listSystemCommand(goos string) (string, error) {
	switch goos {
	case "windows":
		return "'set'", nil
	case "linux", "darwin":
		return "'printenv'", nil
	default:
		return "", unsupportedOSError(goos)
	}
}

func listSysPathCommand(goos string) (string, error) {
	switch goos {
	case "windows":
		return "'echo %PATH%'", nil
	case "linux", "darwin":
		return "'echo $PATH'", nil
	default:
		return "", unsupportedOSError(goos)
	}
}

func unsupportedOSError(goos string) error {
	return fmt.Errorf("%w: %q", ErrUnsupportedOS, goos)
}
