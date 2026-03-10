package listdom

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/command"

	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/pkg/utils"
)

type Handler struct {
	app     *app.Config
	execute *executedom.Handler
}

func NewHandler(app *app.Config) *Handler {

	execute := executedom.NewHandler(app)

	return &Handler{
		app:     app,
		execute: execute,
	}
}

// List all commands in order
func (h *Handler) ListCommands() {
	h.app.Logger.Info("------------------------------------------------")
	print(h.app.OrderKeys, h.app.AllCommandsByName, h.app)
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
	switch runtime.GOOS {
	case "windows":
		return h.execute.PerformExecute(fmt.Sprintf("dir '%s'", folderPath))
	case "linux":
		return h.execute.PerformExecute(fmt.Sprintf("ls '%s'", folderPath))
	}

	return "", nil
}

// List all system variable
func (h *Handler) ListSystem() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return h.execute.PerformExecute("'set'")
	case "linux":
		return h.execute.PerformExecute("'printenv'")
	}

	return "", nil
}

// List system path variable
func (h *Handler) ListSysPath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return h.execute.PerformExecute("'echo %PATH%'")
	case "linux":
		return h.execute.PerformExecute("'echo $PATH'")
	}

	return "", nil
}
