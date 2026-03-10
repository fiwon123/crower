package opendom

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/app"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
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

func (c *Core) Open(args []string) {

	paths := []string{}
	for _, arg := range args {
		matches, err := filepath.Glob(arg)
		if err != nil {
			continue
		}

		if len(matches) > 0 {
			paths = append(paths, matches...)
		} else {
			paths = append(paths, arg)
		}
	}

	c.handler.Open(paths)

	c.app.History.Add(operationsdata.Open, newOpenNote(args))
	historyhelper.Save(c.app)
}

func (c *Core) OpenFile(args []string) {
	currentPath := c.app.CfgFilePath
	if len(args) == 0 {
		args = append(args, currentPath)
	}

	c.handler.Open(args)

	c.app.History.Add(operationsdata.Open, newOpenFolderNote(args))
	historyhelper.Save(c.app)
}

func (c *Core) OpenFolder(args []string) {
	currentPath := c.app.CfgFilePath
	if len(args) == 0 {
		args = append(args, filepath.Dir(currentPath))
	}

	c.handler.Open(args)

	c.app.History.Add(operationsdata.Open, newOpenFolderNote(args))
	historyhelper.Save(c.app)
}

func (c *Core) OpenSystem() {
	err := c.handler.OpenSystem()
	if err != nil {
		c.app.Logger.Error("failed to open system variable: ", "error", err)
		return
	}

	c.app.History.Add(operationsdata.Open, newOpenSystemNote())
	historyhelper.Save(c.app)
}
