package deletedom

import (
	"github.com/fiwon123/crower/internal/app"

	"github.com/fiwon123/crower/internal/helper"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
)

type listHandler interface {
	ListCommands()
}

type Input struct {
	app  *app.Config
	list listHandler
}

func newInput(app *app.Config, list listHandler) *Input {
	return &Input{
		app:  app,
		list: list,
	}
}

// Verify parameters to process delete operation
func (i *Input) CheckDeleteInput(key *string) bool {

	if *key == "" {
		i.list.ListCommands()
		input := helper.GetUserInput("Select Row", helper.IsValidInputKey, i.app).(string)
		*key = input
	}

	var command *command.Data
	if *key != "" {
		command = i.app.AllCommandsByName.Get(*key)
		if command == nil {
			command = i.app.AllCommandsByAlias.Get(*key)
		}
	}

	if command == nil {
		i.list.ListCommands()
		i.app.Logger.Info("Command not found, try to select one.")
		input := helper.GetUserInput("Select Row", helper.IsValidInputKey, i.app).(string)
		*key = input

		command = i.app.AllCommandsByName.Get(*key)
	}

	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("Name:    ", "name", command.Name)
	i.app.Logger.Info("Aliases: ", "alias", command.AllAlias)
	i.app.Logger.Info("Exec:    ", "exec", command.Exec)
	i.app.Logger.Info("")

	ok := helper.GetUserConfirmation("Continue to delete", i.app)
	return ok

}

// Verify parameters to process delete history content operation
func (i *Input) CheckDeleteHistoryContentInput() (history.Content, bool) {

	var content history.Content
	i.app.Logger.Info(i.app.History.GetList())
	content = helper.GetUserInput("Select Row", helper.IsValidContentKey, i.app).(history.Content)

	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("Content")
	i.app.Logger.Info("Version:    ", "version", content.Version)
	i.app.Logger.Info("File:    ", "file", content.File)
	i.app.Logger.Info("Note:    ", "note", content.Note)
	i.app.Logger.Info("")

	ok := helper.GetUserConfirmation("Continue to restore", i.app)
	return content, ok
}
