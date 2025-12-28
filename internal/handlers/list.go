package handlers

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/pkg/utils"
)

func ListCommands(app *app.Data) {
	fmt.Println("------------------------------------------------")
	print(app.OrderKeys, app.AllCommandsByName)
}

func print(orderKeys []string, allCommands command.MapData) {
	fmt.Printf("%-3s %-12s %-16s %-8s \n", "Row", "Name", "Aliases", "Exec")
	fmt.Println("------------------------------------------------")

	for i, key := range orderKeys {
		command := allCommands.Get(key)
		fmt.Printf("%-3d %-12s %-16v %-8s \n", i, command.Name, strings.Join(command.AllAlias, ","), command.Exec)
	}
}

// List all ListProcess running on user operational system (OS).
func ListProcess(args []string, app *app.Data) {

	partName := ""
	if len(args) > 0 {
		partName = args[0]
	}

	err := utils.ListAllProcess(partName, true)
	if err != nil {
		app.LoggerInfo.Error("Error getting processes:", err)
	}
}

func ListHistory(app *app.Data) error {

	app.History.List()

	return nil
}

func ListFolder(folderPath string, app *app.Data) ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute(fmt.Sprintf("'dir %s'", folderPath))
	case "linux":
		return PerformExecute(fmt.Sprintf("'ls %s'", folderPath))
	}

	return nil, nil
}

func ListSystem(*app.Data) ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute("'set'")
	case "linux":
		return PerformExecute("'printenv'")
	}

	return nil, nil
}

func ListSysPath(app *app.Data) ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute("'echo %PATH%'")
	case "linux":
		return PerformExecute("'echo $PATH'")
	}

	return nil, nil
}
