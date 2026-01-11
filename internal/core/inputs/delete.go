package inputs

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/handlers"
)

// Verify parameters to process delete operation
func CheckDeleteInput(key *string, app *app.Data) bool {

	if *key == "" {
		handlers.ListCommands(app)
		input := getUserInput("Select Row", isValidInputKey, app).(string)
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
		input := getUserInput("Select Row", isValidInputKey, app).(string)
		*key = input

		command = app.AllCommandsByName.Get(*key)
	}

	app.Logger.Info("-----------------------------------------")
	app.Logger.Info("Name:    ", "name", command.Name)
	app.Logger.Info("Aliases: ", "alias", command.AllAlias)
	app.Logger.Info("Exec:    ", "exec", command.Exec)
	app.Logger.Info("")

	ok := getUserConfirmation("Continue to delete", app)
	return ok

}

// Verify parameters to process delete history content operation
func CheckDeleteHistoryContentInput(app *app.Data) (history.Content, bool) {

	var content history.Content
	app.Logger.Info(app.History.GetList())
	content = getUserInput("Select Row", isValidContentKey, app).(history.Content)

	app.Logger.Info("-----------------------------------------")
	app.Logger.Info("Content")
	app.Logger.Info("Version:    ", "version", content.Version)
	app.Logger.Info("File:    ", "file", content.File)
	app.Logger.Info("Note:    ", "note", content.Note)
	app.Logger.Info("")

	ok := getUserConfirmation("Continue to restore", app)
	return content, ok
}
