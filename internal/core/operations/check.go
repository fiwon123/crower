package operations

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func CheckNewVersion(currentVersion string, app *app.Data) {
	newVersion, err := handlers.CheckNewVersion(currentVersion, app)
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
}
