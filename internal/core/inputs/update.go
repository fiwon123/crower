package inputs

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

// Verify parameters to process update operation
func CheckUpdateInput(key *string, name *string, allAlias *[]string, exec *string, app *app.Data) bool {

	if *key == "" {
		handlers.ListCommands(app)
		input := getUserInput("Select Row", isValidInputKey, app).(string)
		*key = input
	}

	app.Logger.Info("-----------------------------------------")
	updateCommand := app.AllCommandsByName.Get(*key)
	app.Logger.Info("Name:    ", "name", updateCommand.Name)
	app.Logger.Info("Aliases: ", "alias", updateCommand.AllAlias)
	app.Logger.Info("Exec:    ", "exec", updateCommand.Exec)
	app.Logger.Info("")

	if *name == "" {
		ok := getUserConfirmation("Do you want to update name", app)

		if ok {
			*name = inputName(app)
		}
	}

	if len(*allAlias) == 0 {
		ok := getUserConfirmation("Do you want to update alias", app)

		if ok {
			*allAlias = inputAlias(app)
		}
	}

	if *exec == "" {
		ok := getUserConfirmation("Do you want to update exec", app)

		if ok {
			*exec = inputExec(app)
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

	ok := getUserConfirmation("Continue to update", app)
	return ok
}
