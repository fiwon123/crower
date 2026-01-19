package deleteoperations

import (
	deleteinputs "github.com/fiwon123/crower/internal/core/inputs/delete"
	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"

	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func Delete(args []string, app *app.Data) {

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

	app.History.Add(state.Delete, notes.GenerateDeleteCommandNote(command, args))
	history.Save(app)
}

func performDeleteCommand(key string, app *app.Data) *command.Data {

	command, ok := handlers.DeleteCommand(key, app)
	if !ok {
		app.Logger.Error("Error delete command: ", key)
		return nil
	}

	app.Logger.Info("deleted command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	return command
}

func DeleteLast(op state.MainOperationEnum, app *app.Data) {
	content := history.GetLast(op, app)

	if content == nil {
		crerrors.PrintCommandNotFoundError(app)
		return
	}

	command := performDeleteCommand(content.CommandName, app)
	if command != nil {
		return
	}

	app.History.Add(state.Delete, notes.GenerateDeleteLastNote(op, command))
	history.Save(app)
}

func DeleteSystemVariable(args []string, app *app.Data) {
	newVar := ""
	if len(args) >= 1 {
		newVar = args[0]
	} else {
		crerrors.PrintNotArgs("var name", app)
		return
	}

	out, err := handlers.DeleteSystemVariable(newVar, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	app.Logger.Info(out)

	app.History.Add(state.Delete, notes.GenerateDeleteSystemVariable(args))
	history.Save(app)
}

func DeleteSystemPathVariable(args []string, app *app.Data) {
	newPath := ""
	if len(args) > 0 {
		newPath = args[0]
	} else {
		crerrors.PrintNotArgs("path", app)
		return
	}

	out, err := handlers.DeleteSystemPathVariable(newPath, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	app.Logger.Info(out)

	app.History.Add(state.Delete, notes.GenerateDeleteSystemPathVariable(args))
	history.Save(app)
}

func DeleteFile(args []string, app *app.Data) {
	filePath := ""
	if len(args) > 0 {
		filePath = args[0]
	} else {
		crerrors.PrintNotFileAndOutputPath(app)
		return
	}

	handlers.DeleteFile(filePath, app)

	app.History.Add(state.Execute, notes.GenerateDeleteFileNote(args))
	history.Save(app)
}

func DeleteFolder(args []string, app *app.Data) {
	folderPath := ""
	if len(args) > 0 {
		folderPath = args[0]
	} else {
		crerrors.PrintNotFileAndOutputPath(app)
		return
	}

	handlers.DeleteFolder(folderPath, app)

	app.History.Add(state.Execute, notes.GenerateDeleteFolderNote(args))
	history.Save(app)
}

func DeleteHistoryContent(args []string, app *app.Data) {
	content, ok := deleteinputs.CheckDeleteHistoryContentInput(app)
	if !ok {
		app.Logger.Info("Cancelling Delete History Content...")
		return
	}

	out, err := handlers.DeleteHistoryContent(content, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	app.Logger.Info(out)

	history.SaveOnlyHistory(app)
}
