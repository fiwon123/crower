package extractdom

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/crowererrors"
	extractnotesdata "github.com/fiwon123/crower/internal/data/notes/extract"
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

func (c *Core) Extract(args []string, outDir string) {

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
		crowererrors.PrintEmptyPaths(c.app)
		return
	}

	c.handler.Extract(paths, outDir)

	c.app.History.Add(operationsdata.Execute, extractnotesdata.NewExtractNote(outDir, args))
	historyhelper.Save(c.app)
}
