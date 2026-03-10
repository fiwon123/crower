package updatedom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/crowererrors"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	updatenotesdata "github.com/fiwon123/crower/internal/data/notes/update"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
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

	c.app.History.Add(operationsdata.Update, updatenotesdata.NewUpdateCommmandNote(args, oldCommand, newCommand))
	historyhelper.Save(c.app)
}

func (c *Core) performUpdateCommand(key string, name string, allAlias []string, exec string) (*commanddata.Data, *commanddata.Data) {
	oldCommand, newCommand, err := c.handler.UpdateCommand(key, name, allAlias, exec)
	if err != nil {
		c.app.Logger.Error("Error update command: ", "error", err, "key", key, "name", name, "alias", allAlias, "exec", exec)
		return nil, nil
	}

	c.app.Logger.Info("updated command: ", c.app.AllCommandsByName)
	crowerutils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath)

	return oldCommand, newCommand
}

func (c *Core) UpdateLast(op operationsdata.MainOperationEnum, name string, allAlias []string, exec string) {
	content := historyhelper.GetLast(op, c.app)

	if content == nil {
		crowererrors.PrintCommandNotFoundError(c.app)
		return
	}

	key := content.CommandName

	oldCommand, newCommand := c.performUpdateCommand(key, name, allAlias, exec)
	if oldCommand == nil || newCommand == nil {
		return
	}

	c.app.History.Add(operationsdata.Update, updatenotesdata.NewUpdateLastNote(op, oldCommand, newCommand))
	historyhelper.Save(c.app)
}
