package updateinputs

import (
	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

// Verify parameters to process update operation
func CheckUpdateInput(key *string, name *string, allAlias *[]string, exec *string, app *app.Data) bool {

	if *key == "" {
		handlers.ListCommands(app)
		input := inputs.GetUserInput("Select Row", inputs.IsValidInputKey, app).(string)
		*key = input
	}

	app.Logger.Info("-----------------------------------------")
	updateCommand := app.AllCommandsByName.Get(*key)
	app.Logger.Info("Name:    ", "name", updateCommand.Name)
	app.Logger.Info("Aliases: ", "alias", updateCommand.AllAlias)
	app.Logger.Info("Exec:    ", "exec", updateCommand.Exec)
	app.Logger.Info("")

	if *name == "" {
		ok := inputs.GetUserConfirmation("Do you want to update name", app)

		if ok {
			*name = inputs.InputName(app)
		}
	}

	if len(*allAlias) == 0 {
		ok := inputs.GetUserConfirmation("Do you want to update alias", app)

		if ok {
			*allAlias = inputs.InputAlias(app)
		}
	}

	if *exec == "" {
		ok := inputs.GetUserConfirmation("Do you want to update exec", app)

		if ok {
			*exec = inputs.InputExec(app)
		}
	}
	app.Logger.Info("-----------------------------------------")
	app.Logger.Info("Old Command: ")
	app.Logger.Info("Name:    ", "name", updateCommand.Name)
	app.Logger.Info("Aliases: ", "alias", updateCommand.AllAlias)
	app.Logger.Info("Exec:    ", "exec", updateCommand.Exec)
	app.Logger.Info("-----------------------------------------")
	app.Logger.Info("New Command: ")
	app.Logger.Info("Name:    ", "name", *name)
	app.Logger.Info("Aliases: ", "alias", *allAlias)
	app.Logger.Info("Exec:    ", "exec", *exec)
	app.Logger.Info("")

	ok := inputs.GetUserConfirmation("Continue to update", app)
	return ok
}
