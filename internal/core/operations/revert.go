package operations

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/history"
)

func Revert(app *app.Data) {
	backHistory := app.History.GetBeforeLast()

	if backHistory == nil {
		app.LoggerInfo.Error("Error already in initial history")
		return
	}

	err := history.RevertTo(backHistory, app)
	if err != nil {
		app.LoggerInfo.Error("Error revert history %v", err)
		return
	}
	app.LoggerInfo.Info("reverted to history version ", backHistory.Version)
	history.Save(app)
}
