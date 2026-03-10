package executeinputs

import (
	inputscore "github.com/fiwon123/crower/internal/core/inputs"
	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"

	listhandlers "github.com/fiwon123/crower/internal/handlers/list"
)

// Verify parameters to process execute operation
func CheckExecuteInput(key *string, params *[]string, app *appdata.Data) bool {
	if *key == "" {
		listhandlers.ListCommands(app)
		input := inputscore.GetUserInput("Select Row", inputscore.IsValidInputKey, app).(string)
		*key = input
	}

	var command *commanddata.Data
	if *key != "" {
		command = app.AllCommandsByName.Get(*key)
		if command == nil {
			command = app.AllCommandsByAlias.Get(*key)
		}
	}

	if command == nil {
		listhandlers.ListCommands(app)
		app.Logger.Info("Command not found, try to select one.")
		input := inputscore.GetUserInput("Select Row", inputscore.IsValidInputKey, app).(string)
		*key = input

		command = app.AllCommandsByName.Get(*key)
	}

	app.Logger.Info("-----------------------------------------")
	app.Logger.Info("Name:    ", "name", command.Name)
	app.Logger.Info("Aliases: ", "alias", command.AllAlias)
	app.Logger.Info("Exec:    ", "exec", command.Exec)
	app.Logger.Info("")

	ok := inputscore.GetUserConfirmation("Continue to execute", app)
	return ok
}
