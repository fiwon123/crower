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

// CheckNewVersion reports whether a release newer than currentVersion exists.
func (c *Core) CheckNewVersion(currentVersion string) error {
	newVersion, err := c.handler.CheckNewVersion(currentVersion)
	if err != nil {
		return err
	}

	c.app.Logger.Info("Current Version: ", currentVersion)

	if newVersion == currentVersion {
		c.app.Logger.Info("already up-to-date")
		return nil
	}

	c.app.Logger.Info("New Version Found: ", newVersion)
	c.app.Logger.Info("Check: https://github.com/fiwon123/crower/releases/latest")

	c.app.History.Add(operations.Check, newCheckNote())

	return c.app.Save()
}
