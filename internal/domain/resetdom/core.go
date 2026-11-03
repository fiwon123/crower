package resetdom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/operation"

	"github.com/fiwon123/crower/internal/interfaces"
)

type Core struct {
	app     *app.Config
	handler *Handler
	input   *Input
}

func NewCore(app *app.Config, list interfaces.ListHandler) *Core {

	handler := NewHandler(app)
	input := NewInput(app, list)

	return &Core{
		handler: handler,
		app:     app,
		input:   input,
	}
}

// Reset erases every managed command.
func (c *Core) Reset() error {
	ok, err := c.input.CheckResetInput()
	if err != nil {
		return err
	}

	if !ok {
		c.app.Logger.Info("Cancelling reset...")
		return nil
	}

	c.app.Logger.Info("reset all commands")
	c.handler.Reset()

	if err := c.app.SaveCommands(); err != nil {
		return err
	}

	c.app.History.Add(operation.Reset, "", newResetNote())

	return c.app.Save()
}
