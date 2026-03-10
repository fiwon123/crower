package updateoperations

import (
	updateinputs "github.com/fiwon123/crower/internal/core/inputs/update"
	"github.com/fiwon123/crower/internal/crowererrors"
	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	updatenotesdata "github.com/fiwon123/crower/internal/data/notes/update"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	updatehandlers "github.com/fiwon123/crower/internal/handlers/update"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

func UpdateCommand(args []string, name string, allAlias []string, exec string, app *appdata.Data) {

	key := ""
	if len(args) != 0 {
		key = args[0]
	}

	ok := updateinputs.CheckUpdateInput(&key, &name, &allAlias, &exec, app)
	if !ok {
		app.Logger.Info("Cancelling update...")
		return
	}

	oldCommand, newCommand := performUpdateCommand(key, name, allAlias, exec, app)
	if oldCommand == nil || newCommand == nil {
		return
	}

	app.History.Add(operationsdata.Update, updatenotesdata.NewUpdateCommmandNote(args, oldCommand, newCommand))
	historyhelper.Save(app)
}

func performUpdateCommand(key string, name string, allAlias []string, exec string, app *appdata.Data) (*commanddata.Data, *commanddata.Data) {
	oldCommand, newCommand, err := updatehandlers.UpdateCommand(key, name, allAlias, exec, app)
	if err != nil {
		app.Logger.Error("Error update command: ", "error", err, "key", key, "name", name, "alias", allAlias, "exec", exec)
		return nil, nil
	}

	app.Logger.Info("updated command: ", app.AllCommandsByName)
	crowerutils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	return oldCommand, newCommand
}

func UpdateLast(op operationsdata.MainOperationEnum, name string, allAlias []string, exec string, app *appdata.Data) {
	content := historyhelper.GetLast(op, app)

	if content == nil {
		crowererrors.PrintCommandNotFoundError(app)
		return
	}

	key := content.CommandName

	oldCommand, newCommand := performUpdateCommand(key, name, allAlias, exec, app)
	if oldCommand == nil || newCommand == nil {
		return
	}

	app.History.Add(operationsdata.Update, updatenotesdata.NewUpdateLastNote(op, oldCommand, newCommand))
	historyhelper.Save(app)
}
