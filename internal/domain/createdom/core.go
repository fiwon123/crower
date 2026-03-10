package createdom

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/crowererrors"
	"github.com/fiwon123/crower/internal/helper"

	commanddata "github.com/fiwon123/crower/internal/data/command"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	openhandlers "github.com/fiwon123/crower/internal/handlers/open"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

type CreateCore struct {
	app     *app.Config
	handler *CreateHandler
}

func NewCore(app *app.Config) *CreateCore {

	createHandler := NewHandler(app)

	return &CreateCore{
		handler: createHandler,
		app:     app,
	}
}

// Verify parameters to process create operation
func (c *CreateCore) CheckCreateInput(name *string, alias *[]string, exec *string) {
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

func (c *CreateCore) CreateCommand(allAlias []string, args []string) {
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

	c.app.History.Add(operationsdata.Create, NewCreateCommandNote(command, args))
	historyhelper.Save(c.app)
}

func (c *CreateCore) performCreateCommand(name string, allAlias []string, exec string) *commanddata.Data {
	command, err := c.handler.CreateCommand(name, allAlias, exec)

	if err != nil {
		c.app.Logger.Error("Error add command: ", "error", err, "name", name, "alias", allAlias, "exec", exec)
		return nil
	}

	crowerutils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath)
	c.app.Logger.Info("added new command: ", "allCommands", c.app.AllCommandsByName)

	return command
}

func (c *CreateCore) CreateProcess(name string, args []string) {
	command, err := c.handler.CreateProcess(name, args)
	if err != nil {
		c.app.Logger.Error("Error add command by process: ", "error", err, "name", name, "args", args)
		return
	}

	crowerutils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath)
	c.app.Logger.Info("added new command by process: ", "allCommands", c.app.AllCommandsByName)

	c.app.History.Add(operationsdata.Create, NewCreateProcessNote(command, args))
	historyhelper.Save(c.app)
}

func (c *CreateCore) CreateSystemVariable(args []string) {
	newVar := ""
	value := ""
	if len(args) >= 2 {
		newVar = args[0]
		value = args[1]
	} else {
		crowererrors.PrintNotArgs("var name and var value", c.app)
		return
	}

	out, err := c.handler.CreateSystemVariable(newVar, value)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operationsdata.Create, NewCreateSystemVariableNote(args))
	historyhelper.Save(c.app)
}

func (c *CreateCore) CreateSystemPathVariable(args []string) {
	newPath := ""
	if len(args) > 0 {
		newPath = args[0]
	} else {
		crowererrors.PrintNotArgs("path", c.app)
		return
	}

	out, err := c.handler.CreateSystemPathVariable(newPath)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operationsdata.Create, NewCreateSystemPathVariableNote(args))
	historyhelper.Save(c.app)
}

func (c *CreateCore) CreateFile(args []string) {
	for _, path := range args {
		err := c.handler.CreateFile(path)
		if err != nil {
			c.app.Logger.Error(err.Error())
		}
	}

	c.app.History.Add(operationsdata.Create, GenerateCreateFile(args))
	historyhelper.Save(c.app)
}

func (c *CreateCore) CreateFolder(args []string) {
	for _, path := range args {
		err := c.handler.CreateFolder(path)
		if err != nil {
			c.app.Logger.Error(err.Error())
		}
	}

	c.app.History.Add(operationsdata.Create, NewCreateFolder(args))
	historyhelper.Save(c.app)
}

func (c *CreateCore) CreateLastCommand(op operationsdata.MainOperationEnum, args []string) {
	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		crowererrors.PrintNotArgs("name", c.app)
		return
	}

	content := historyhelper.GetLast(op, c.app)

	if content == nil {
		crowererrors.PrintCommandNotFoundError(c.app)
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

	c.app.History.Add(operationsdata.Create, NewCreateCommandLastExecuteNote(command))
	historyhelper.Save(c.app)
}

func (c *CreateCore) CreateScriptCommand(args []string) {
	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		crowererrors.PrintNotArgs("name", c.app)
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

	openhandlers.Open([]string{filepath.Dir(scriptFilePath)}, c.app)

	c.app.History.Add(operationsdata.Create, NewCreateScriptCommandNote(command, args))
	historyhelper.Save(c.app)
}
