package updateinputs

import (
	inputscore "github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	listhandlers "github.com/fiwon123/crower/internal/handlers/list"
)

// Verify parameters to process update operation
func CheckUpdateInput(key *string, name *string, allAlias *[]string, exec *string, app *app.Data) bool {

	if *key == "" {
		listhandlers.ListCommands(app)
		input := inputscore.GetUserInput("Select Row", inputscore.IsValidInputKey, app).(string)
		*key = input
	}

	app.Logger.Info("-----------------------------------------")
	updateCommand := app.AllCommandsByName.Get(*key)
	app.Logger.Info("Name:    ", "name", updateCommand.Name)
	app.Logger.Info("Aliases: ", "alias", updateCommand.AllAlias)
	app.Logger.Info("Exec:    ", "exec", updateCommand.Exec)
	app.Logger.Info("")

	if *name == "" {
		ok := inputscore.GetUserConfirmation("Do you want to update name", app)

		if ok {
			*name = inputscore.InputName(app)
		}
	}

	if len(*allAlias) == 0 {
		ok := inputscore.GetUserConfirmation("Do you want to update alias", app)

		if ok {
			*allAlias = inputscore.InputAlias(app)
		}
	}

	if *exec == "" {
		ok := inputscore.GetUserConfirmation("Do you want to update exec", app)

		if ok {
			*exec = inputscore.InputExec(app)
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

	ok := inputscore.GetUserConfirmation("Continue to update", app)
	return ok
}
