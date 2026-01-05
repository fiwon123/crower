package operations

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
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
		crerrors.PrintEmptyPaths()
		return
	}

	handlers.Extract(paths, outDir, app)
}
