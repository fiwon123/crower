package updatedom

import (
	"github.com/fiwon123/crower/internal/app"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/operations"

	"github.com/fiwon123/crower/internal/interfaces"
	"github.com/fiwon123/crower/pkg/utils"
)

type Core struct {
	app     *app.Config
	handler *Handler
	input   *Input
}

func NewCore(app *app.Config, listHandler interfaces.ListHandler) *Core {

	handler := NewHandler(app)
	input := newInput(app, listHandler)

	return &Core{
		handler: handler,
		app:     app,
		input:   input,
	}
}

// UpdateCommand changes the name, aliases or exec of a managed command.
func (c *Core) UpdateCommand(args []string, name string, allAlias []string, exec string) error {

	key := ""
	if len(args) != 0 {
		key = args[0]
	}

	ok, err := c.input.CheckUpdateInput(&key, &name, &allAlias, &exec)
	if err != nil {
		return err
	}

	if !ok {
		c.app.Logger.Info("Cancelling update...")
		return nil
	}

	oldCommand, newCommand, err := c.performUpdateCommand(key, name, allAlias, exec)
	if err != nil {
		return err
	}

	c.app.History.Add(operations.Update, newUpdateCommmandNote(args, oldCommand, newCommand))

	return c.app.Save()
}

func (c *Core) performUpdateCommand(
	key string,
	name string,
	allAlias []string,
	exec string,
) (*command.Data, *command.Data, error) {
	oldCommand, newCommand, err := c.handler.UpdateCommand(key, name, allAlias, exec)
	if err != nil {
		return nil, nil, err
	}

	if err := utils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath); err != nil {
		return nil, nil, err
	}

	c.app.Logger.Info("updated command: ", key)

	return oldCommand, newCommand, nil
}

// UpdateLast updates the command touched by the last operation of the given
// kind.
func (c *Core) UpdateLast(op operations.MainOperationEnum, name string, allAlias []string, exec string) error {
	content := c.app.GetLast(op)

	if content == nil {
		return commandNotFound(op)
	}

	key := content.CommandName

	oldCommand, newCommand, err := c.performUpdateCommand(key, name, allAlias, exec)
	if err != nil {
		return err
	}

	c.app.History.Add(operations.Update, newUpdateLastNote(op, oldCommand, newCommand))

	return c.app.Save()
}
