package executeinputs

import (
	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/handlers"
)

// Verify parameters to process execute operation
func CheckExecuteInput(key *string, params *[]string, app *app.Data) bool {
	if *key == "" {
		handlers.ListCommands(app)
		input := inputs.GetUserInput("Select Row", inputs.IsValidInputKey, app).(string)
		*key = input
	}

	var command *command.Data
	if *key != "" {
		command = app.AllCommandsByName.Get(*key)
		if command == nil {
			command = app.AllCommandsByAlias.Get(*key)
		}
	}

	if command == nil {
		handlers.ListCommands(app)
		app.Logger.Info("Command not found, try to select one.")
		input := inputs.GetUserInput("Select Row", inputs.IsValidInputKey, app).(string)
		*key = input

		command = app.AllCommandsByName.Get(*key)
	}

	app.Logger.Info("-----------------------------------------")
	app.Logger.Info("Name:    ", "name", command.Name)
	app.Logger.Info("Aliases: ", "alias", command.AllAlias)
	app.Logger.Info("Exec:    ", "exec", command.Exec)
	app.Logger.Info("")

	ok := inputs.GetUserConfirmation("Continue to execute", app)
	return ok
}
