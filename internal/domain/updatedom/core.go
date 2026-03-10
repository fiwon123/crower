package updatedom

import (
	"github.com/fiwon123/crower/internal/app"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/errors"

	"github.com/fiwon123/crower/internal/interfaces"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

type Core struct {
	app     *app.Config
	handler *Handler
	input   *Input
}

func NewCore(app *app.Config, listHandler interfaces.ListHandler) *Core {

	handler := NewHandler(app)
	input := NewInput(app, listHandler)

	return &Core{
		handler: handler,
		app:     app,
		input:   input,
	}
}

func (c *Core) UpdateCommand(args []string, name string, allAlias []string, exec string) {

	key := ""
	if len(args) != 0 {
		key = args[0]
	}

	ok := c.input.CheckUpdateInput(&key, &name, &allAlias, &exec)
	if !ok {
		c.app.Logger.Info("Cancelling update...")
		return
	}

	oldCommand, newCommand := c.performUpdateCommand(key, name, allAlias, exec)
	if oldCommand == nil || newCommand == nil {
		return
	}

	c.app.History.Add(operations.Update, newUpdateCommmandNote(args, oldCommand, newCommand))
	c.app.Save()
}

func (c *Core) performUpdateCommand(key string, name string, allAlias []string, exec string) (*command.Data, *command.Data) {
	oldCommand, newCommand, err := c.handler.UpdateCommand(key, name, allAlias, exec)
	if err != nil {
		c.app.Logger.Error("Error update command: ", "error", err, "key", key, "name", name, "alias", allAlias, "exec", exec)
		return nil, nil
	}

	c.app.Logger.Info("updated command: ", c.app.AllCommandsByName)
	crowerutils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath)

	return oldCommand, newCommand
}

func (c *Core) UpdateLast(op operations.MainOperationEnum, name string, allAlias []string, exec string) {
	content := c.app.GetLast(op)

	if content == nil {
		errors.PrintCommandNotFoundError(c.app)
		return
	}

	key := content.CommandName

	oldCommand, newCommand := c.performUpdateCommand(key, name, allAlias, exec)
	if oldCommand == nil || newCommand == nil {
		return
	}

	c.app.History.Add(operations.Update, newUpdateLastNote(op, oldCommand, newCommand))
	c.app.Save()
}
