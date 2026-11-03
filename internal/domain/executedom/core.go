package executedom

import (
	"strings"

	"github.com/fiwon123/crower/internal/app"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/errors"

	"github.com/fiwon123/crower/internal/helper"
)

type listHandler interface {
	CommandsTable() string
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

// ExecuteCommand runs a managed command.
func (c *Core) ExecuteCommand(args []string) error {
	var params []string

	key := ""
	if len(args) > 0 {
		key = args[0]
		params = args[1:]
	} else {
		ok, err := c.CheckExecuteInput(&key, &params)
		if err != nil {
			return err
		}

		if !ok {
			c.app.Logger.Info("Cancelling execute...")
			return nil
		}
	}

	output, executed, err := c.handler.Execute(key, params)
	if err != nil {
		return err
	}

	c.report(output)

	c.app.History.Add(operations.Execute, newExecuteCommandNote(executed))

	return c.app.Save()
}

// ExecuteLast runs the command touched by the last operation of the given kind.
func (c *Core) ExecuteLast(op operations.MainOperationEnum, args []string) error {
	content := c.app.GetLast(op)

	if content == nil {
		return errors.CommandNotFound(string(op))
	}

	output, executed, err := c.handler.Execute(content.CommandName, args)
	if err != nil {
		return err
	}

	c.report(output)

	c.app.History.Add(operations.Execute, NewExecuteLastNote(op, executed))

	return c.app.Save()
}

// report writes a command's own output to the result stream.
//
// The output belongs to the executed command, so it is data a caller may pipe,
// not a diagnostic about the run. Printing it through the diagnostic writer put
// it on stderr and appended a space.
func (c *Core) report(output string) {
	trimmed := strings.TrimRight(output, "\n")
	if trimmed == "" {
		return
	}

	c.app.Logger.Out(trimmed)
}
