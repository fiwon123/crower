package executedom

import (
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/errors"
)

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
