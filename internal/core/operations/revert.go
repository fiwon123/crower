package operations

import (
	"fmt"
	"strconv"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/history"
)

func Revert(args []string, app *app.Data) {

	steps := 1
	var err error
	if len(args) > 0 {
		steps, err = strconv.Atoi(args[0])

		if err != nil {
			crerrors.PrintNotArgs("steps int number")
			return
		}

	} else {
		crerrors.PrintNotArgs("steps int number")
		return
	}

	ok, err := inputs.CheckRevertInput(steps, app)
	if !ok {
		fmt.Printf("error %v\n", err)
		return
	}

	backHistory, err := app.History.GetBeforeLast(steps)

	if err != nil {
		app.LoggerInfo.Error("error: ", err)
		return
	}

	err = history.RevertTo(backHistory, app)
	if err != nil {
		app.LoggerInfo.Error("Error revert history %v", err)
		return
	}
	app.LoggerInfo.Info("reverted to history version ", backHistory.Version)
	history.Save(app)
}
