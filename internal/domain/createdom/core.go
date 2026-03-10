package createdom

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/opendom"
	"github.com/fiwon123/crower/internal/errors"

	"github.com/fiwon123/crower/internal/helper"

	commanddata "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

type Core struct {
	app         *app.Config
	handler     *Handler
	openHandler opendom.Handler
}

func NewCore(app *app.Config, openHandler opendom.Handler) *Core {

	handler := NewHandler(app)

	return &Core{
		handler:     handler,
		app:         app,
		openHandler: openHandler,
	}
}

// Verify parameters to process create operation
func (c *Core) CheckCreateInput(name *string, alias *[]string, exec *string) {
	if *name == "" {
		*name = helper.GetUserInput("New Name ", helper.IsValidInput, c.app).(string)
	}

	if len(*alias) == 0 {
		ok := helper.GetUserConfirmation("Do you want to add alias", c.app)

		if ok {
			*alias = helper.InputAlias(c.app)
		}
	}

	if *exec == "" {
		*exec = helper.GetUserInput("New Exec ", helper.IsValidInput, c.app).(string)
	}

}

func (c *Core) CreateCommand(allAlias []string, args []string) {
	name := ""
	exec := ""
	if len(args) == 2 {
		name = args[0]
		exec = args[1]
	} else {
		c.CheckCreateInput(&name, &allAlias, &exec)
	}

	command := c.performCreateCommand(name, allAlias, exec)
	if command == nil {
		return
	}

	c.app.History.Add(operations.Create, newCreateCommandNote(command, args))
	c.app.Save()
}

func (c *Core) performCreateCommand(name string, allAlias []string, exec string) *commanddata.Data {
	command, err := c.handler.CreateCommand(name, allAlias, exec)

	if err != nil {
		c.app.Logger.Error("Error add command: ", "error", err, "name", name, "alias", allAlias, "exec", exec)
		return nil
	}

	crowerutils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath)
	c.app.Logger.Info("added new command: ", "allCommands", c.app.AllCommandsByName)

	return command
}

func (c *Core) CreateProcess(name string, args []string) {
	command, err := c.handler.CreateProcess(name, args)
	if err != nil {
		c.app.Logger.Error("Error add command by process: ", "error", err, "name", name, "args", args)
		return
	}

	crowerutils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath)
	c.app.Logger.Info("added new command by process: ", "allCommands", c.app.AllCommandsByName)

	c.app.History.Add(operations.Create, NewCreateProcessNote(command, args))
	c.app.Save()
}

func (c *Core) CreateSystemVariable(args []string) {
	newVar := ""
	value := ""
	if len(args) >= 2 {
		newVar = args[0]
		value = args[1]
	} else {
		errors.PrintNotArgs("var name and var value", c.app)
		return
	}

	out, err := c.handler.CreateSystemVariable(newVar, value)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operations.Create, NewCreateSystemVariableNote(args))
	c.app.Save()
}

func (c *Core) CreateSystemPathVariable(args []string) {
	newPath := ""
	if len(args) > 0 {
		newPath = args[0]
	} else {
		errors.PrintNotArgs("path", c.app)
		return
	}

	out, err := c.handler.CreateSystemPathVariable(newPath)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operations.Create, NewCreateSystemPathVariableNote(args))
	c.app.Save()
}

func (c *Core) CreateFile(args []string) {
	for _, path := range args {
		err := c.handler.CreateFile(path)
		if err != nil {
			c.app.Logger.Error(err.Error())
		}
	}

	c.app.History.Add(operations.Create, GenerateCreateFile(args))
	c.app.Save()
}

func (c *Core) CreateFolder(args []string) {
	for _, path := range args {
		err := c.handler.CreateFolder(path)
		if err != nil {
			c.app.Logger.Error(err.Error())
		}
	}

	c.app.History.Add(operations.Create, NewCreateFolder(args))
	c.app.Save()
}

func (c *Core) CreateLastCommand(op operations.MainOperationEnum, args []string) {
	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		errors.PrintNotArgs("name", c.app)
		return
	}

	content := c.app.GetLast(op)

	if content == nil {
		errors.PrintCommandNotFoundError(c.app)
		return
	}

	exec := ""
	key := content.CommandName
	if key == "" {
		splitted := strings.SplitSeq(content.Note, ";")
		for keyValRaw := range splitted {
			keyVal := strings.Split(keyValRaw, "=")
			if keyVal[0] == "exec" {
				exec = keyVal[1]
				break
			}
		}
	}

	command := c.performCreateCommand(name, []string{}, exec)
	if command == nil {
		return
	}

	c.app.History.Add(operations.Create, NewCreateCommandLastExecuteNote(command))
	c.app.Save()
}

func (c *Core) CreateScriptCommand(args []string) {
	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		errors.PrintNotArgs("name", c.app)
	}

	scriptFilePath, err := c.handler.CreateScriptCommand(name)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	var command *commanddata.Data
	switch runtime.GOOS {
	case "windows":
		command = c.performCreateCommand(name, []string{}, scriptFilePath)
	case "linux":
		command = c.performCreateCommand(name, []string{}, scriptFilePath)
	}

	if command == nil {
		return
	}

	c.openHandler.Open([]string{filepath.Dir(scriptFilePath)})

	c.app.History.Add(operations.Create, NewCreateScriptCommandNote(command, args))
	c.app.Save()
}
