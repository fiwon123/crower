package updatedom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/helper"
	"github.com/fiwon123/crower/internal/interfaces"
)

type Input struct {
	app         *app.Config
	listHandler interfaces.ListHandler
}

func newInput(app *app.Config, listHandler interfaces.ListHandler) *Input {
	return &Input{
		app:         app,
		listHandler: listHandler,
	}
}

// Verify parameters to process update operation
func (i *Input) CheckUpdateInput(key *string, name *string, allAlias *[]string, exec *string) bool {

	if *key == "" {
		i.listHandler.ListCommands()
		input := helper.GetUserInput("Select Row", helper.IsValidInputKey, i.app).(string)
		*key = input
	}

	i.app.Logger.Info("-----------------------------------------")
	updateCommand := i.app.AllCommandsByName.Get(*key)
	i.app.Logger.Info("Name:    ", "name", updateCommand.Name)
	i.app.Logger.Info("Aliases: ", "alias", updateCommand.AllAlias)
	i.app.Logger.Info("Exec:    ", "exec", updateCommand.Exec)
	i.app.Logger.Info("")

	if *name == "" {
		ok := helper.GetUserConfirmation("Do you want to update name", i.app)

		if ok {
			*name = helper.InputName(i.app)
		}
	}

	if len(*allAlias) == 0 {
		ok := helper.GetUserConfirmation("Do you want to update alias", i.app)

		if ok {
			*allAlias = helper.InputAlias(i.app)
		}
	}

	if *exec == "" {
		ok := helper.GetUserConfirmation("Do you want to update exec", i.app)

		if ok {
			*exec = helper.InputExec(i.app)
		}
	}
	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("Old Command: ")
	i.app.Logger.Info("Name:    ", "name", updateCommand.Name)
	i.app.Logger.Info("Aliases: ", "alias", updateCommand.AllAlias)
	i.app.Logger.Info("Exec:    ", "exec", updateCommand.Exec)
	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("New Command: ")
	i.app.Logger.Info("Name:    ", "name", *name)
	i.app.Logger.Info("Aliases: ", "alias", *allAlias)
	i.app.Logger.Info("Exec:    ", "exec", *exec)
	i.app.Logger.Info("")

	ok := helper.GetUserConfirmation("Continue to update", i.app)
	return ok
}
