package extractoperations

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/crowererrors"
	appdata "github.com/fiwon123/crower/internal/data/app"
	extractnotesdata "github.com/fiwon123/crower/internal/data/notes/extract"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	extracthandlers "github.com/fiwon123/crower/internal/handlers/extract"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
)

func Extract(args []string, outDir string, app *appdata.Data) {

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

	app.History.Add(operationsdata.Execute, extractnotesdata.NewExtractNote(outDir, args))
	historyhelper.Save(app)
}
