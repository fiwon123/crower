package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func CreateCommand(payload payload.Data, app *app.Data) {
	inputs.CheckCreateInput(&payload.Name, &payload.Alias, &payload.Exec, app)

	command, err := handlers.CreateCommand(payload.Name, payload.Alias, payload.Exec, payload.Args, app)

	if err != nil {
		app.LoggerInfo.Error("Error add command: ", err, payload)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command: ", app.AllCommandsByName)

	app.History.Add(operation.Create, command.Name, notes.GenerateAddNote(command))
	history.Save(app)
}

func CreateProcess(payload payload.Data, app *app.Data) {
	command, err := handlers.CreateProcess(payload.Name, payload.Args, app)
	if err != nil {
		app.LoggerInfo.Error("Error add command by process: ", err, payload)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command by process: ", app.AllCommandsByName)

	app.History.Add(operation.Create, command.Name, notes.GenerateAddProcessNote(command))
	history.Save(app)
}

func CreateFile(args []string, app *app.Data) {
	currentPath := "./"
	fileName := ""
	if len(args) > 0 {

		if len(args) > 1 {
			currentPath = args[0]
			fileName = args[1]
		} else {
			fileName = args[0]
		}
	} else {
		fmt.Println("file name and/or folder path not specified")
		return
	}

	handlers.CreateFile(currentPath, fileName, app)
}

func CreateFolder(args []string, app *app.Data) {
	currentPath := "./"
	folderName := ""
	if len(args) > 0 {
		if len(args) > 1 {
			currentPath = args[0]
			folderName = args[1]
		} else {
			folderName = args[0]
		}
	} else {
		fmt.Println("file name and/or folder path not specified")
		return
	}

	handlers.CreateFolder(currentPath, folderName, app)
}
