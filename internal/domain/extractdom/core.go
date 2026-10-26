package extractdom

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/errors"
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

// Extract expands every archive in args into outDir, or beside each archive
// when outDir is empty.
func (c *Core) Extract(args []string, outDir string) error {

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

	if len(paths) == 0 {
		return errors.EmptyPaths()
	}

	if err := c.handler.Extract(paths, outDir); err != nil {
		return err
	}

	c.app.History.Add(operations.Execute, newExtractNote(outDir, args))

	return c.app.Save()
}
