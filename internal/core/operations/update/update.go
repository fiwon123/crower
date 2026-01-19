package updateoperations

import (
	updateinputs "github.com/fiwon123/crower/internal/core/inputs/update"
	"github.com/fiwon123/crower/internal/crowererrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/state"
	updatehandlers "github.com/fiwon123/crower/internal/handlers/update"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

func UpdateCommand(args []string, name string, allAlias []string, exec string, app *app.Data) {

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

	app.History.Add(state.Update, notes.GenerateUpdateCommmandNote(args, oldCommand, newCommand))
	history.Save(app)
}

func performUpdateCommand(key string, name string, allAlias []string, exec string, app *app.Data) (*command.Data, *command.Data) {
	oldCommand, newCommand, err := updatehandlers.UpdateCommand(key, name, allAlias, exec, app)
	if err != nil {
		app.Logger.Error("Error update command: ", "error", err, "key", key, "name", name, "alias", allAlias, "exec", exec)
		return nil, nil
	}

	app.Logger.Info("updated command: ", app.AllCommandsByName)
	crowerutils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	return oldCommand, newCommand
}

func UpdateLast(op state.MainOperationEnum, name string, allAlias []string, exec string, app *app.Data) {
	content := history.GetLast(op, app)

	if content == nil {
		crowererrors.PrintCommandNotFoundError(app)
		return
	}

	key := content.CommandName

	oldCommand, newCommand := performUpdateCommand(key, name, allAlias, exec, app)
	if oldCommand == nil || newCommand == nil {
		return
	}

	app.History.Add(state.Update, notes.GenerateUpdateLastNote(op, oldCommand, newCommand))
	history.Save(app)
}
