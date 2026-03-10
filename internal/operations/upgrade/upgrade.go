package upgradeoperations

import (
	"fmt"

	appdata "github.com/fiwon123/crower/internal/data/app"
	upgradenotesdata "github.com/fiwon123/crower/internal/data/notes/upgrade"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	checkhandlers "github.com/fiwon123/crower/internal/handlers/check"
	upgradehandlers "github.com/fiwon123/crower/internal/handlers/upgrade"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
)

func UpgradeApp(currentVersion string, app *appdata.Data) {
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

	app.History.Add(operationsdata.Upgrade, upgradenotesdata.NewUpgradeNote())
	historyhelper.Save(app)
}
