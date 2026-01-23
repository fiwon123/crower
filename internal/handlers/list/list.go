package listhandlers

import (
	"fmt"
	"runtime"
	"strings"

	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	executehandlers "github.com/fiwon123/crower/internal/handlers/execute"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

// List all commands in order
func ListCommands(app *appdata.Data) {
	app.Logger.Info("------------------------------------------------")
	print(app.OrderKeys, app.AllCommandsByName, app)
}

func print(orderKeys []string, allCommands commanddata.MapData, app *appdata.Data) {
	app.Logger.Info(fmt.Sprintf("%-3s %-12s %-16s %-8s \n", "Row", "Name", "Aliases", "Exec"))
	app.Logger.Info("------------------------------------------------")

	for i, key := range orderKeys {
		command := allCommands.Get(key)
		app.Logger.Info(fmt.Sprintf("%-3d %-12s %-16v %-8s \n", i, command.Name, strings.Join(command.AllAlias, ","), command.Exec))
	}
}

// List all ListProcess running on user operational system (OS).
func ListProcess(args []string, app *appdata.Data) error {

	partName := ""
	if len(args) > 0 {
		partName = args[0]
	}

	out, err := crowerutils.GetAllProcess(partName, true)
	if err != nil {
		app.Logger.Error("Error getting processes:", err)
		return err
	}

	app.Logger.Info(out)

	return nil

}

// List all history
func ListHistory(app *appdata.Data) error {

	app.Logger.Info(app.History.GetList())

	return nil
}

// List all files and folder from a folderpath
func ListFolder(folderPath string, app *appdata.Data) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return executehandlers.PerformExecute(fmt.Sprintf("dir '%s'", folderPath), app)
	case "linux":
		return executehandlers.PerformExecute(fmt.Sprintf("ls '%s'", folderPath), app)
	}

	return "", nil
}

// List all system variable
func ListSystem(app *appdata.Data) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return executehandlers.PerformExecute("'set'", app)
	case "linux":
		return executehandlers.PerformExecute("'printenv'", app)
	}

	return "", nil
}

// List system path variable
func ListSysPath(app *appdata.Data) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return executehandlers.PerformExecute("'echo %PATH%'", app)
	case "linux":
		return executehandlers.PerformExecute("'echo $PATH'", app)
	}

	return "", nil
}
