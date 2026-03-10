package executedom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/crowererrors"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"

	"github.com/fiwon123/crower/internal/helper"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
)

type listHandler interface {
	ListCommands()
}

type Core struct {
	app     *app.Config
	handler *Handler
	list    listHandler
}

func NewCore(app *app.Config, list listHandler) *Core {

	handler := NewHandler(app)

	return &Core{
		handler: handler,
		app:     app,
		list:    list,
	}
}

// Verify parameters to process execute operation
func (c *Core) CheckExecuteInput(key *string, params *[]string) bool {
	if *key == "" {
		c.list.ListCommands()
		input := helper.GetUserInput("Select Row", helper.IsValidInputKey, c.app).(string)
		*key = input
	}

	var command *commanddata.Data
	if *key != "" {
		command = c.app.AllCommandsByName.Get(*key)
		if command == nil {
			command = c.app.AllCommandsByAlias.Get(*key)
		}
	}

	if command == nil {
		c.list.ListCommands()
		c.app.Logger.Info("Command not found, try to select one.")
		input := helper.GetUserInput("Select Row", helper.IsValidInputKey, c.app).(string)
		*key = input

		command = c.app.AllCommandsByName.Get(*key)
	}

	c.app.Logger.Info("-----------------------------------------")
	c.app.Logger.Info("Name:    ", "name", command.Name)
	c.app.Logger.Info("Aliases: ", "alias", command.AllAlias)
	c.app.Logger.Info("Exec:    ", "exec", command.Exec)
	c.app.Logger.Info("")

	ok := helper.GetUserConfirmation("Continue to execute", c.app)
	return ok
}

func (c *Core) ExecuteCommand(args []string) {
	var params []string
	key := ""
	if len(args) > 0 {
		c.app.Logger.Debug("args", args)
		key = args[0]
		params = args[1:]
	} else {
		ok := c.CheckExecuteInput(&key, &params)
		if !ok {
			c.app.Logger.Info("Cancelling execute...")
			return
		}
	}

	output, command, err := c.handler.Execute(key, params)
	c.assertExecute(output, command, err)

	c.app.History.Add(operationsdata.Execute, newExecuteCommandNote(command))
	historyhelper.Save(c.app)
}

func (c *Core) ExecuteLast(op operationsdata.MainOperationEnum, args []string) {
	content := historyhelper.GetLast(op, c.app)

	if content == nil {
		crowererrors.PrintCommandNotFoundError(c.app)
		return
	}

	output, command, err := c.handler.Execute(content.CommandName, args)
	c.assertExecute(output, command, err)

	c.app.History.Add(operationsdata.Execute, NewExecuteLastNote(op, command))
	historyhelper.Save(c.app)
}

func (c *Core) assertExecute(output string, command *commanddata.Data, err error) {
	if err != nil {
		c.app.Logger.Error("Error trying to run command: ", "out", string(output), "err", err)
		return
	}
	c.app.Logger.Info(string(output))
}
