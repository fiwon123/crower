package deletedom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
)

type listHandler interface {
	CommandsTable() string
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

// CheckDeleteInput resolves the command to delete and asks for confirmation.
//
// It returns false with a nil error when the user declines. An error means no
// command could be chosen -- an unknown key with no way to ask, or no input
// available -- which the command layer reports as a failure.
func (i *Input) CheckDeleteInput(key *string) (bool, error) {

	if *key == "" {
		if err := i.selectRow(key); err != nil {
			return false, err
		}
	}

	selected := i.find(*key)

	if selected == nil {
		i.app.Logger.Info("Command not found, try to select one.")

		if err := i.selectRow(key); err != nil {
			return false, err
		}

		selected = i.find(*key)
		if selected == nil {
			return false, errors.CommandNotFound(*key)
		}
	}

	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("Name:    ", selected.Name)
	i.app.Logger.Info("Aliases: ", selected.AllAlias)
	i.app.Logger.Info("Exec:    ", selected.Exec)
	i.app.Logger.Info("")

	return helper.GetUserConfirmation("Continue to delete", i.app)
}

// selectRow asks for a row of the printed command table.
func (i *Input) selectRow(key *string) error {
	i.app.Logger.Info(i.list.CommandsTable())

	value, err := helper.GetUserInput("Select Row", helper.IsValidInputKey, i.app)
	if err != nil {
		return err
	}

	selected, ok := value.(string)
	if !ok || selected == "" {
		return fmt.Errorf("selected row holds no command")
	}

	*key = selected
	return nil
}

// find looks a command up by name and then by alias.
func (i *Input) find(key string) *command.Data {
	if key == "" {
		return nil
	}

	if found := i.app.AllCommandsByName.Get(key); found != nil {
		return found
	}

	return i.app.AllCommandsByAlias.Get(key)
}

// CheckDeleteHistoryContentInput asks which history entry to drop.
func (i *Input) CheckDeleteHistoryContentInput() (history.Content, bool, error) {
	var content history.Content

	i.app.Logger.Info(i.app.History.GetList())

	value, err := helper.GetUserInput("Select Row", helper.IsValidContentKey, i.app)
	if err != nil {
		return content, false, err
	}

	selected, ok := value.(history.Content)
	if !ok {
		return content, false, fmt.Errorf("selected row is not a history entry")
	}
	content = selected

	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("Content")
	i.app.Logger.Info("Version:    ", content.Version)
	i.app.Logger.Info("File:    ", content.File)
	i.app.Logger.Info("Note:    ", content.Note)
	i.app.Logger.Info("")

	confirmed, err := helper.GetUserConfirmation("Continue to delete history content", i.app)
	return content, confirmed, err
}
