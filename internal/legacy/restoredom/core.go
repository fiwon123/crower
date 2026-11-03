package restoredom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/operation"
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

// Restore brings back the command stored in a history entry.
func (c *Core) Restore(args []string) error {
	key := ""
	if len(args) > 0 {
		key = args[0]
	}

	content, ok, err := c.input.CheckRestoreInput()
	if err != nil {
		return err
	}

	if !ok {
		c.app.Logger.Info("Cancelling restore...")
		return nil
	}

	out, err := c.handler.RestoreHistory(key, content)
	if err != nil {
		return err
	}

	// The restored command has to reach the command file, not only the map this
	// process holds: Save() writes the history snapshot and index, and nothing
	// else writes the configuration.
	if err := c.app.SaveCommands(); err != nil {
		return err
	}

	c.app.Logger.Info("restored command: ", out)

	c.app.History.Add(operation.Restore, "", newRestoreNote(args))

	return c.app.Save()
}
