package upgradedom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/app"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/domain/checkdom"

	historyhelper "github.com/fiwon123/crower/internal/helper/history"
)

type Core struct {
	app          *app.Config
	handler      *Handler
	checkHandler *checkdom.Handler
}

func NewCore(app *app.Config) *Core {

	handler := NewHandler(app)
	checkHandler := checkdom.NewHandler(app)

	return &Core{
		app:          app,
		handler:      handler,
		checkHandler: checkHandler,
	}
}

func (c *Core) UpgradeApp(currentVersion string) {
	newVersion, err := c.checkHandler.CheckNewVersion(currentVersion)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	if newVersion == currentVersion {
		c.app.Logger.Info("Current Version: ", "currentVersion", currentVersion)
		c.app.Logger.Info("already up-to-date")
		return
	}

	err = c.handler.UpgradeApp(newVersion)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	c.app.Logger.Info(fmt.Sprintf("crower upgraded from %s to %s \n", currentVersion, newVersion))

	c.app.History.Add(operationsdata.Upgrade, newUpgradeNote())
	historyhelper.Save(c.app)
}
