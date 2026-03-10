package deletedom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/crowererrors"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	deletenotesdata "github.com/fiwon123/crower/internal/data/notes/delete"

	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

type Core struct {
	app     *app.Config
	handler *Handler
	input   *Input
}

func NewCore(app *app.Config) *Core {

	handler := NewHandler(app)
	input := NewInput(app)

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

	c.app.History.Add(operationsdata.Delete, deletenotesdata.NewDeleteCommandNote(command, args))
	historyhelper.Save(c.app)
}

func (c *Core) performDeleteCommand(key string) *commanddata.Data {

	command, ok := c.handler.DeleteCommand(key)
	if !ok {
		c.app.Logger.Error("Error delete command: ", key)
		return nil
	}

	c.app.Logger.Info("deleted command: ", c.app.AllCommandsByName)
	crowerutils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath)

	return command
}

func (c *Core) DeleteLast(op operationsdata.MainOperationEnum) {
	content := historyhelper.GetLast(op, c.app)

	if content == nil {
		crowererrors.PrintCommandNotFoundError(c.app)
		return
	}

	command := c.performDeleteCommand(content.CommandName)
	if command != nil {
		return
	}

	c.app.History.Add(operationsdata.Delete, deletenotesdata.NewDeleteLastNote(op, command))
	historyhelper.Save(c.app)
}

func (c *Core) DeleteSystemVariable(args []string) {
	newVar := ""
	if len(args) >= 1 {
		newVar = args[0]
	} else {
		crowererrors.PrintNotArgs("var name", c.app)
		return
	}

	out, err := c.handler.DeleteSystemVariable(newVar)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operationsdata.Delete, deletenotesdata.NewDeleteSystemVariable(args))
	historyhelper.Save(c.app)
}

func (c *Core) DeleteSystemPathVariable(args []string) {
	newPath := ""
	if len(args) > 0 {
		newPath = args[0]
	} else {
		crowererrors.PrintNotArgs("path", c.app)
		return
	}

	out, err := c.handler.DeleteSystemPathVariable(newPath)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operationsdata.Delete, deletenotesdata.NewDeleteSystemPathVariable(args))
	historyhelper.Save(c.app)
}

func (c *Core) DeleteFile(args []string) {
	filePath := ""
	if len(args) > 0 {
		filePath = args[0]
	} else {
		crowererrors.PrintNotFileAndOutputPath(c.app)
		return
	}

	c.handler.DeleteFile(filePath)

	c.app.History.Add(operationsdata.Execute, deletenotesdata.NewDeleteFileNote(args))
	historyhelper.Save(c.app)
}

func (c *Core) DeleteFolder(args []string) {
	folderPath := ""
	if len(args) > 0 {
		folderPath = args[0]
	} else {
		crowererrors.PrintNotFileAndOutputPath(c.app)
		return
	}

	c.handler.DeleteFolder(folderPath)

	c.app.History.Add(operationsdata.Execute, deletenotesdata.GenerateDeleteFolderNote(args))
	historyhelper.Save(c.app)
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

	historyhelper.SaveOnlyHistory(c.app)
}
