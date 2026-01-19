package extractoperations

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/crowererrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	extracthandlers "github.com/fiwon123/crower/internal/handlers/extract"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func Extract(args []string, outDir string, app *app.Data) {

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
		crowererrors.PrintEmptyPaths(app)
		return
	}

	extracthandlers.Extract(paths, outDir, app)

	app.History.Add(state.Execute, notes.GenerateExtractNote(outDir, args))
	history.Save(app)
}
