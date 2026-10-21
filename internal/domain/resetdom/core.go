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

func (c *Core) Reset() {
	ok := c.input.CheckResetInput()

	if !ok {
		c.app.Logger.Info("Cancelling reset...")
		return
	}

	c.app.Logger.Info("reset all commands: ", c.app.AllCommandsByName)
	c.handler.Reset()
	utils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath)

	c.app.History.Add(operations.Reset, newResetNote())
	c.app.Save()
}
