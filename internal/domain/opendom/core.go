package opendom

import (
	"path/filepath"

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

// Open opens every path with the platform's default application.
func (c *Core) Open(args []string) error {

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

	if err := c.handler.Open(paths); err != nil {
		return err
	}

	c.app.History.Add(operations.Open, newOpenNote(args))

	return c.app.Save()
}

// OpenFile opens a file, defaulting to the configuration file.
func (c *Core) OpenFile(args []string) error {
	currentPath := c.app.CfgFilePath
	if len(args) == 0 {
		args = append(args, currentPath)
	}

	if err := c.handler.Open(args); err != nil {
		return err
	}

	c.app.History.Add(operations.Open, newOpenFolderNote(args))

	return c.app.Save()
}

// OpenFolder opens a folder, defaulting to the configuration folder.
func (c *Core) OpenFolder(args []string) error {
	currentPath := c.app.CfgFilePath
	if len(args) == 0 {
		args = append(args, filepath.Dir(currentPath))
	}

	if err := c.handler.Open(args); err != nil {
		return err
	}

	c.app.History.Add(operations.Open, newOpenFolderNote(args))

	return c.app.Save()
}

// OpenSystem opens the platform's system variable editor.
func (c *Core) OpenSystem() error {
	if err := c.handler.OpenSystem(); err != nil {
		return err
	}

	c.app.History.Add(operations.Open, newOpenSystemNote())

	return c.app.Save()
}
