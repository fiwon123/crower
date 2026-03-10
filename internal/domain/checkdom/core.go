package checkdom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
)

type Core struct {
	app     *app.Config
	handler *Handler
}

func NewCore(app *app.Config) *Core {

	handler := NewHandler(app)

	return &Core{
		handler: handler,
		app:     app,
	}
}

func (c *Core) CheckNewVersion(currentVersion string) {
	newVersion, err := c.handler.CheckNewVersion(currentVersion)
	if err != nil {
		c.app.Logger.Error(err.Error())
		return
	}

	if newVersion == currentVersion {
		c.app.Logger.Info("Current Version: ", "currentVersion", currentVersion)
		c.app.Logger.Info("already up-to-date")
		return
	}

	c.app.Logger.Info("Current Version: ", "currentVersion", currentVersion)
	c.app.Logger.Info("New Version Found: ", "newVersion", newVersion)
	c.app.Logger.Info("Check: https://github.com/fiwon123/crower/releases/latest")

	c.app.History.Add(operations.Check, newCheckNote())
	c.app.Save()
}
