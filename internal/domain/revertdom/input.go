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

// Verify parameters to process revert operation
func (i *Input) CheckRevertInput(steps int) (bool, error) {
	if steps < 0 {
		return false, fmt.Errorf("steps is negative number")
	}

	stopIndex := i.app.History.GetIndexFromLastTo(steps)
	if stopIndex < 0 {
		return false, fmt.Errorf("can't revert steps is greater than quantity of history registry")
	}

	i.app.Logger.Info("")
	i.app.Logger.Info("Deleted History")
	i.app.Logger.Info(i.app.History.GetListLastHistory(steps))

	i.app.Logger.Info("")
	i.app.Logger.Info("New History")
	i.app.Logger.Info(i.app.History.GetListFirstHistory(stopIndex))

	i.app.Logger.Info("")
	i.app.Logger.Info("History will revert ", "steps", steps)
	ok := helper.GetUserConfirmation("Continue to revert", i.app)

	if !ok {
		return false, fmt.Errorf("cancelling revert...")
	}

	return ok, nil
}
