package handlers

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/pkg/utils"
)

// List all commands in order
func ListCommands(app *app.Data) {
	app.Logger.Info("------------------------------------------------")
	print(app.OrderKeys, app.AllCommandsByName, app)
}

func print(orderKeys []string, allCommands command.MapData, app *app.Data) {
	app.Logger.Info(fmt.Sprintf("%-3s %-12s %-16s %-8s \n", "Row", "Name", "Aliases", "Exec"))
	app.Logger.Info("------------------------------------------------")

	for i, key := range orderKeys {
		command := allCommands.Get(key)
		app.Logger.Info(fmt.Sprintf("%-3d %-12s %-16v %-8s \n", i, command.Name, strings.Join(command.AllAlias, ","), command.Exec))
	}
}

// List all ListProcess running on user operational system (OS).
func ListProcess(args []string, app *app.Data) error {

	partName := ""
	if len(args) > 0 {
		partName = args[0]
	}

	err := utils.ListAllProcess(partName, true)
	if err != nil {
		app.Logger.Error("Error getting processes:", err)
		return err
	}

	return nil

}

// List all history
func ListHistory(app *app.Data) error {

	app.Logger.Info(app.History.GetList())

	return nil
}

// List all files and folder from a folderpath
func ListFolder(folderPath string, app *app.Data) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute(fmt.Sprintf("dir '%s'", folderPath), app)
	case "linux":
		return PerformExecute(fmt.Sprintf("ls '%s'", folderPath), app)
	}

	return "", nil
}

// List all system variable
func ListSystem(app *app.Data) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute("'set'", app)
	case "linux":
		return PerformExecute("'printenv'", app)
	}

	return "", nil
}

// List system path variable
func ListSysPath(app *app.Data) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute("'echo %PATH%'", app)
	case "linux":
		return PerformExecute("'echo $PATH'", app)
	}

	return "", nil
}
