package checkoperations

import (
	appdata "github.com/fiwon123/crower/internal/data/app"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	checkhandlers "github.com/fiwon123/crower/internal/handlers/check"

	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func CheckNewVersion(currentVersion string, app *appdata.Data) {
	newVersion, err := checkhandlers.CheckNewVersion(currentVersion, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	if newVersion == currentVersion {
		app.Logger.Info("Current Version: ", "currentVersion", currentVersion)
		app.Logger.Info("already up-to-date")
		return
	}

	app.Logger.Info("Current Version: ", "currentVersion", currentVersion)
	app.Logger.Info("New Version Found: ", "newVersion", newVersion)
	app.Logger.Info("Check: https://github.com/fiwon123/crower/releases/latest")

	app.History.Add(operationsdata.Check, notes.GenerateCheckNote())
	history.Save(app)
}
