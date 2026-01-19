package upgradeoperations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	checkhandlers "github.com/fiwon123/crower/internal/handlers/check"
	upgradehandlers "github.com/fiwon123/crower/internal/handlers/upgrade"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func UpgradeApp(currentVersion string, app *app.Data) {
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

	err = upgradehandlers.UpgradeApp(newVersion, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	app.Logger.Info(fmt.Sprintf("crower upgraded from %s to %s \n", currentVersion, newVersion))

	app.History.Add(state.Upgrade, notes.GenerateUpgradeNote())
	history.Save(app)
}
