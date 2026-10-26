package revertdom

import (
	"fmt"
	"strconv"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/errors"
)

type Core struct {
	app   *app.Config
	input *Input
}

func NewCore(app *app.Config) *Core {

	input := NewInput(app)

	return &Core{
		app:   app,
		input: input,
	}
}

// Revert moves the configuration back by steps history entries.
func (c *Core) Revert(args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("steps")
	}

	steps, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("steps must be a number: %q", args[0])
	}

	ok, err := c.input.CheckRevertInput(steps)
	if err != nil {
		return err
	}

	if !ok {
		c.app.Logger.Info("Cancelling revert...")
		return nil
	}

	backHistory, err := c.app.History.GetBeforeLast(steps)
	if err != nil {
		return err
	}

	if err := c.app.RevertTo(backHistory); err != nil {
		return fmt.Errorf("revert history: %w", err)
	}

	c.app.Logger.Info("reverted to history version ", backHistory.Version)

	c.app.History.Add(operations.Revert, newRevertNote(args))

	return c.app.Save()
}
