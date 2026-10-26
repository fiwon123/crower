package upgradedom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/domain/checkdom"
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

// UpgradeApp replaces the running binary with the newest release.
func (c *Core) UpgradeApp(currentVersion string) error {
	newVersion, err := c.checkHandler.CheckNewVersion(currentVersion)
	if err != nil {
		return err
	}

	c.app.Logger.Info("Current Version: ", currentVersion)

	if newVersion == currentVersion {
		c.app.Logger.Info("already up-to-date")
		return nil
	}

	if err := c.handler.UpgradeApp(newVersion); err != nil {
		return err
	}

	c.app.Logger.Info("crower upgraded from ", currentVersion, " to ", newVersion)

	c.app.History.Add(operations.Upgrade, newUpgradeNote())

	return c.app.Save()
}
