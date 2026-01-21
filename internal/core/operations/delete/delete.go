package deleteoperations

import (
	deleteinputs "github.com/fiwon123/crower/internal/core/inputs/delete"
	"github.com/fiwon123/crower/internal/crowererrors"
	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	deletehandlers "github.com/fiwon123/crower/internal/handlers/delete"

	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

func Delete(args []string, app *appdata.Data) {

	key := ""
	if len(args) > 0 {
		key = args[0]
	}

	ok := deleteinputs.CheckDeleteInput(&key, app)
	if !ok {
		app.Logger.Info("Cancelling delete...")
		return
	}

	command := performDeleteCommand(key, app)
	if command != nil {
		return
	}

	app.History.Add(operationsdata.Delete, notes.GenerateDeleteCommandNote(command, args))
	history.Save(app)
}

func performDeleteCommand(key string, app *appdata.Data) *commanddata.Data {

	command, ok := deletehandlers.DeleteCommand(key, app)
	if !ok {
		app.Logger.Error("Error delete command: ", key)
		return nil
	}

	app.Logger.Info("deleted command: ", app.AllCommandsByName)
	crowerutils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	return command
}

func DeleteLast(op operationsdata.MainOperationEnum, app *appdata.Data) {
	content := history.GetLast(op, app)

	if content == nil {
		crowererrors.PrintCommandNotFoundError(app)
		return
	}

	command := performDeleteCommand(content.CommandName, app)
	if command != nil {
		return
	}

	app.History.Add(operationsdata.Delete, notes.GenerateDeleteLastNote(op, command))
	history.Save(app)
}

func DeleteSystemVariable(args []string, app *appdata.Data) {
	newVar := ""
	if len(args) >= 1 {
		newVar = args[0]
	} else {
		crowererrors.PrintNotArgs("var name", app)
		return
	}

	out, err := deletehandlers.DeleteSystemVariable(newVar, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	app.Logger.Info(out)

	app.History.Add(operationsdata.Delete, notes.GenerateDeleteSystemVariable(args))
	history.Save(app)
}

func DeleteSystemPathVariable(args []string, app *appdata.Data) {
	newPath := ""
	if len(args) > 0 {
		newPath = args[0]
	} else {
		crowererrors.PrintNotArgs("path", app)
		return
	}

	out, err := deletehandlers.DeleteSystemPathVariable(newPath, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	app.Logger.Info(out)

	app.History.Add(operationsdata.Delete, notes.GenerateDeleteSystemPathVariable(args))
	history.Save(app)
}

func DeleteFile(args []string, app *appdata.Data) {
	filePath := ""
	if len(args) > 0 {
		filePath = args[0]
	} else {
		crowererrors.PrintNotFileAndOutputPath(app)
		return
	}

	deletehandlers.DeleteFile(filePath, app)

	app.History.Add(operationsdata.Execute, notes.GenerateDeleteFileNote(args))
	history.Save(app)
}

func DeleteFolder(args []string, app *appdata.Data) {
	folderPath := ""
	if len(args) > 0 {
		folderPath = args[0]
	} else {
		crowererrors.PrintNotFileAndOutputPath(app)
		return
	}

	deletehandlers.DeleteFolder(folderPath, app)

	app.History.Add(operationsdata.Execute, notes.GenerateDeleteFolderNote(args))
	history.Save(app)
}

func DeleteHistoryContent(args []string, app *appdata.Data) {
	content, ok := deleteinputs.CheckDeleteHistoryContentInput(app)
	if !ok {
		app.Logger.Info("Cancelling Delete History Content...")
		return
	}

	out, err := deletehandlers.DeleteHistoryContent(content, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	app.Logger.Info(out)

	history.SaveOnlyHistory(app)
}
