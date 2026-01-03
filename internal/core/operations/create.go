package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func CreateCommand(name string, allAlias []string, exec string, args []string, app *app.Data) {
	inputs.CheckCreateInput(&name, &allAlias, &exec, app)

	command, err := handlers.CreateCommand(name, allAlias, exec, args, app)

	if err != nil {
		app.LoggerInfo.Error("Error add command: ", err, name, allAlias, exec, args)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command: ", app.AllCommandsByName)

	app.History.Add(operation.Create, command.Name, notes.GenerateAddNote(command))
	history.Save(app)
}

func CreateProcess(name string, args []string, app *app.Data) {
	command, err := handlers.CreateProcess(name, args, app)
	if err != nil {
		app.LoggerInfo.Error("Error add command by process: ", err, name, args)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command by process: ", app.AllCommandsByName)

	app.History.Add(operation.Create, command.Name, notes.GenerateAddProcessNote(command))
	history.Save(app)
}

func CreateFile(args []string, app *app.Data) {
	filePath := ""
	if len(args) > 0 {

		filePath = args[0]
	} else {
		fmt.Println("file name and/or folder path not specified")
		return
	}

	handlers.CreateFile(filePath, app)
}

func CreateFolder(args []string, app *app.Data) {
	folderPath := ""
	if len(args) > 0 {
		folderPath = args[0]
	} else {
		fmt.Println("file name and/or folder path not specified")
		return
	}

	handlers.CreateFolder(folderPath, app)
}
