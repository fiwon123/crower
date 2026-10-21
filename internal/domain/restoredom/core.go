package restoredom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
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
		handler: handler,
		app:     app,
		input:   input,
	}
}

func (c *Core) Restore(args []string) {
	key := ""
	if len(args) > 0 {
		key = args[0]
	}

	content, ok := c.input.CheckRestoreInput()
	if !ok {
		c.app.Logger.Info("Cancelling Restore...")
		return
	}

	out, err := c.handler.RestoreHistory(key, content)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info("restored command: ", "out", out)

	c.app.History.Add(operations.Restore, newRestoreNote(out))
	c.app.Save()
}
