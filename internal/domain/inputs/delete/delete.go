package deleteinputs

import (
	coreinputs "github.com/fiwon123/crower/internal/core/inputs"
	inputscore "github.com/fiwon123/crower/internal/core/inputs"
	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"

	historydata "github.com/fiwon123/crower/internal/data/history"
	listhandlers "github.com/fiwon123/crower/internal/handlers/list"
)

// Verify parameters to process delete operation
func CheckDeleteInput(key *string, app *appdata.Data) bool {

	if *key == "" {
		listhandlers.ListCommands(app)
		input := coreinputs.GetUserInput("Select Row", coreinputs.IsValidInputKey, app).(string)
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
		input := coreinputs.GetUserInput("Select Row", coreinputs.IsValidInputKey, app).(string)
		*key = input

		command = app.AllCommandsByName.Get(*key)
	}

	app.Logger.Info("-----------------------------------------")
	app.Logger.Info("Name:    ", "name", command.Name)
	app.Logger.Info("Aliases: ", "alias", command.AllAlias)
	app.Logger.Info("Exec:    ", "exec", command.Exec)
	app.Logger.Info("")

	ok := coreinputs.GetUserConfirmation("Continue to delete", app)
	return ok

}

// Verify parameters to process delete history content operation
func CheckDeleteHistoryContentInput(app *appdata.Data) (historydata.Content, bool) {

	var content historydata.Content
	app.Logger.Info(app.History.GetList())
	content = inputscore.GetUserInput("Select Row", inputscore.IsValidContentKey, app).(historydata.Content)

	app.Logger.Info("-----------------------------------------")
	app.Logger.Info("Content")
	app.Logger.Info("Version:    ", "version", content.Version)
	app.Logger.Info("File:    ", "file", content.File)
	app.Logger.Info("Note:    ", "note", content.Note)
	app.Logger.Info("")

	ok := inputscore.GetUserConfirmation("Continue to restore", app)
	return content, ok
}
