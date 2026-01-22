package revertoperations

import (
	"strconv"

	revertinputs "github.com/fiwon123/crower/internal/core/inputs/revert"
	"github.com/fiwon123/crower/internal/crowererrors"
	appdata "github.com/fiwon123/crower/internal/data/app"
	revertnotesdata "github.com/fiwon123/crower/internal/data/notes/revert"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
)

func Revert(args []string, app *appdata.Data) {

	steps := 1
	var err error
	if len(args) > 0 {
		steps, err = strconv.Atoi(args[0])

		if err != nil {
			crowererrors.PrintNotArgs("steps int number", app)
			return
		}

	} else {
		crowererrors.PrintNotArgs("steps int number", app)
		return
	}

	ok, err := revertinputs.CheckRevertInput(steps, app)
	if !ok {
		app.Logger.Error(err.Error())
		return
	}

	backHistory, err := app.History.GetBeforeLast(steps)

	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	err = historyhelper.RevertTo(backHistory, app)
	if err != nil {
		app.Logger.Error("Error revert history %v", err)
		return
	}
	app.Logger.Info("reverted to history version ", backHistory.Version)
	historyhelper.Save(app)

	app.History.Add(operationsdata.Revert, revertnotesdata.NewRevertNote(args))
	historyhelper.Save(app)
}
