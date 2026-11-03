package executedom

import (
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"
)

// listHandler is the table a run asks the user to choose from. It is an
// interface so this package does not depend on listdom for its prompt.
type listHandler interface {
	CommandsTable() string
}

// CheckExecuteInput resolves the command to execute and asks for confirmation.
func (c *Core) CheckExecuteInput(key *string, params *[]string) (bool, error) {
	if *key == "" {
		if err := c.selectRow(key); err != nil {
			return false, err
		}
	}

	selected := c.find(*key)

	if selected == nil {
		c.app.Logger.Info("Command not found, try to select one.")

		if err := c.selectRow(key); err != nil {
			return false, err
		}

		selected = c.find(*key)
		if selected == nil {
			return false, errors.CommandNotFound(*key)
		}
	}

	c.app.Logger.Info("-----------------------------------------")
	c.app.Logger.Info("Name:    ", selected.Name)
	c.app.Logger.Info("Aliases: ", selected.AllAlias)
	c.app.Logger.Info("Exec:    ", selected.Exec)
	c.app.Logger.Info("")

	return helper.GetUserConfirmation("Continue to execute", c.app)
}

// selectRow asks for a row of the printed command table.
func (c *Core) selectRow(key *string) error {
	c.app.Logger.Info(c.list.CommandsTable())

	value, err := helper.GetUserInput("Select Row", helper.IsValidInputKey, c.app)
	if err != nil {
		return err
	}

	selected, ok := value.(string)
	if !ok || selected == "" {
		return errors.InvalidRows()
	}

	*key = selected
	return nil
}

// find looks a command up by name and then by alias.
func (c *Core) find(key string) *command.Data {
	if key == "" {
		return nil
	}

	if found := c.app.AllCommandsByName.Get(key); found != nil {
		return found
	}

	return c.app.AllCommandsByAlias.Get(key)
}
