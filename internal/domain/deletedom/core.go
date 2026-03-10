package deletedom

import (
	"github.com/fiwon123/crower/internal/app"

	"github.com/fiwon123/crower/internal/errors"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

type Core struct {
	app     *app.Config
	handler *Handler
	input   *Input
}

func NewCore(app *app.Config, list listHandler) *Core {

	handler := NewHandler(app)
	input := NewInput(app, list)

	return &Core{
		app:     app,
		handler: handler,
		input:   input,
	}
}

func (c *Core) Delete(args []string) {

	key := ""
	if len(args) > 0 {
		key = args[0]
	}

	ok := c.input.CheckDeleteInput(&key)
	if !ok {
		c.app.Logger.Info("Cancelling delete...")
		return
	}

	command := c.performDeleteCommand(key)
	if command != nil {
		return
	}

	c.app.History.Add(operations.Delete, NewDeleteCommandNote(command, args))
	c.app.Save()
}

func (c *Core) performDeleteCommand(key string) *command.Data {

	command, ok := c.handler.DeleteCommand(key)
	if !ok {
		c.app.Logger.Error("Error delete command: ", key)
		return nil
	}

	c.app.Logger.Info("deleted command: ", c.app.AllCommandsByName)
	crowerutils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath)

	return command
}

func (c *Core) DeleteLast(op operations.MainOperationEnum) {
	content := c.app.GetLast(op)

	if content == nil {
		errors.PrintCommandNotFoundError(c.app)
		return
	}

	command := c.performDeleteCommand(content.CommandName)
	if command != nil {
		return
	}

	c.app.History.Add(operations.Delete, NewDeleteLastNote(op, command))
	c.app.Save()
}

func (c *Core) DeleteSystemVariable(args []string) {
	newVar := ""
	if len(args) >= 1 {
		newVar = args[0]
	} else {
		errors.PrintNotArgs("var name", c.app)
		return
	}

	out, err := c.handler.DeleteSystemVariable(newVar)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operations.Delete, NewDeleteSystemVariable(args))
	c.app.Save()
}

func (c *Core) DeleteSystemPathVariable(args []string) {
	newPath := ""
	if len(args) > 0 {
		newPath = args[0]
	} else {
		errors.PrintNotArgs("path", c.app)
		return
	}

	out, err := c.handler.DeleteSystemPathVariable(newPath)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operations.Delete, NewDeleteSystemPathVariable(args))
	c.app.Save()
}

func (c *Core) DeleteFile(args []string) {
	filePath := ""
	if len(args) > 0 {
		filePath = args[0]
	} else {
		errors.PrintNotFileAndOutputPath(c.app)
		return
	}

	c.handler.DeleteFile(filePath)

	c.app.History.Add(operations.Execute, newDeleteFileNote(args))
	c.app.Save()
}

func (c *Core) DeleteFolder(args []string) {
	folderPath := ""
	if len(args) > 0 {
		folderPath = args[0]
	} else {
		errors.PrintNotFileAndOutputPath(c.app)
		return
	}

	c.handler.DeleteFolder(folderPath)

	c.app.History.Add(operations.Execute, generateDeleteFolderNote(args))
	c.app.Save()
}

func (c *Core) DeleteHistoryContent(args []string) {
	content, ok := c.input.CheckDeleteHistoryContentInput()
	if !ok {
		c.app.Logger.Info("Cancelling Delete History Content...")
		return
	}

	out, err := c.handler.DeleteHistoryContent(content)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info(out)

	c.app.SaveOnlyHistory()
}
