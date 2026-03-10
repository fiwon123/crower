package revertdom

import (
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

func (c *Core) Revert(args []string) {

	steps := 1
	var err error
	if len(args) > 0 {
		steps, err = strconv.Atoi(args[0])

		if err != nil {
			errors.PrintNotArgs("steps int number", c.app)
			return
		}

	} else {
		errors.PrintNotArgs("steps int number", c.app)
		return
	}

	ok, err := c.input.CheckRevertInput(steps)
	if !ok {
		c.app.Logger.Error(err.Error())
		return
	}

	backHistory, err := c.app.History.GetBeforeLast(steps)

	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	err = c.app.RevertTo(backHistory)
	if err != nil {
		c.app.Logger.Error("Error revert history %v", err)
		return
	}
	c.app.Logger.Info("reverted to history version ", backHistory.Version)
	c.app.Save()

	c.app.History.Add(operations.Revert, newRevertNote(args))
	c.app.Save()
}
