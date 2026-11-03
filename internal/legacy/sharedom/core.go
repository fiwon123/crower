package sharedom

import (
	"github.com/fiwon123/crower/internal/app"
)

type Core struct {
	app     *app.Config
	handler *Handler
}

func NewCore(app *app.Config) *Core {
	h := NewHandler(app)
	return &Core{
		app:     app,
		handler: h,
	}
}

func (c *Core) Export(name string, outPath string) error {
	return c.handler.Export(name, outPath)
}

func (c *Core) Import(inPath string, overwrite bool) error {
	return c.handler.Import(inPath, overwrite)
}

func (c *Core) List() error {
	return c.handler.ListRegistry()
}

func (c *Core) Pull(name string) error {
	return c.handler.Pull(name)
}
