package revertdom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/helper"
)

type Input struct {
	app *app.Config
}

func NewInput(app *app.Config) *Input {
	return &Input{
		app: app,
	}
}

// CheckRevertInput validates the number of steps and asks for confirmation.
//
// It returns false with a nil error when the user declines: a declined action is
// a choice, not a failure. An error means the revert cannot be attempted at
// all, or that nobody could answer.
func (i *Input) CheckRevertInput(steps int) (bool, error) {
	if steps < 0 {
		return false, fmt.Errorf("steps is negative number")
	}

	stopIndex := i.app.History.GetIndexFromLastTo(steps)
	if stopIndex < 0 {
		return false, fmt.Errorf("can't revert: steps is greater than quantity of history registry")
	}

	i.app.Logger.Info("")
	i.app.Logger.Info("Deleted History")
	i.app.Logger.Info(i.app.History.GetListLastHistory(steps))

	i.app.Logger.Info("")
	i.app.Logger.Info("New History")
	i.app.Logger.Info(i.app.History.GetListFirstHistory(stopIndex))

	i.app.Logger.Info("")
	i.app.Logger.Info("History will revert ", steps)

	return helper.GetUserConfirmation("Continue to revert", i.app)
}
