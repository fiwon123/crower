package revertinputs

import (
	"fmt"

	inputscore "github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
)

// Verify parameters to process revert operation
func CheckRevertInput(steps int, app *app.Data) (bool, error) {
	if steps < 0 {
		return false, fmt.Errorf("steps is negative number")
	}

	stopIndex := app.History.GetIndexFromLastTo(steps)
	if stopIndex < 0 {
		return false, fmt.Errorf("can't revert steps is greater than quantity of history registry")
	}

	app.Logger.Info("")
	app.Logger.Info("Deleted History")
	app.Logger.Info(app.History.GetListLastHistory(steps))

	app.Logger.Info("")
	app.Logger.Info("New History")
	app.Logger.Info(app.History.GetListFirstHistory(stopIndex))

	app.Logger.Info("")
	app.Logger.Info("History will revert ", "steps", steps)
	ok := inputscore.GetUserConfirmation("Continue to revert", app)

	if !ok {
		return false, fmt.Errorf("cancelling revert...")
	}

	return ok, nil
}
