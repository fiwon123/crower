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

func Delete(payload payload.Data, app *app.Data) {

	ok := inputs.CheckDeleteInput(&payload.Name, &payload.Alias, app)
	if !ok {
		fmt.Println("Cancelling delete...")
		return
	}

	key := payload.Name
	if key == "" {
		if len(payload.Alias) > 0 {
			key = payload.Alias[0]
		}
	}

	command, ok := handlers.DeleteCommand(key, app)
	if !ok {
		app.LoggerInfo.Error("Error delete command: ", payload)
		return
	}

	app.LoggerInfo.Info("deleted command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(operation.Delete, command.Name, notes.GenerateDeleteNote(command))
	history.Save(app)
}

func DeleteLastCreate(payload payload.Data, app *app.Data) {
	payload.Op = operation.Create
	content := history.GetLast(payload.Op, app)

	if content == nil {
		fmt.Println("command doesn't exist")
		return
	}

	payload.Name = content.CommandName

	payload.Op = operation.Delete
	Delete(payload, app)
}

func DeleteLastUpdate(payload payload.Data, app *app.Data) {
	payload.Op = operation.Update
	content := history.GetLast(payload.Op, app)

	if content == nil {
		fmt.Println("command doesn't exist")
		return
	}

	payload.Name = content.CommandName

	payload.Op = operation.Delete
	Delete(payload, app)
}

func DeleteLastExecute(payload payload.Data, app *app.Data) {
	payload.Op = operation.Execute
	content := history.GetLast(payload.Op, app)

	if content == nil {
		fmt.Println("command doesn't exist")
		return
	}

	payload.Name = content.CommandName

	payload.Op = operation.Delete
	Delete(payload, app)
}

func DeleteFile(args []string, app *app.Data) {
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

	handlers.DeleteFile(currentPath, fileName, app)
}

func DeleteFolder(args []string, app *app.Data) {
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

	handlers.DeleteFolder(currentPath, folderName, app)
}
