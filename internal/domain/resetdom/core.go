package resetdom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"

	"github.com/fiwon123/crower/internal/interfaces"
	"github.com/fiwon123/crower/pkg/utils"
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

	if err := utils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath); err != nil {
		return err
	}

	c.app.History.Add(operations.Reset, newResetNote())

	return c.app.Save()
}
