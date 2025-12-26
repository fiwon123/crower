package operations

import (
	"path"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func Open(args []string, app *app.Data) {
	currentPath := app.CfgFilePath
	if len(args) > 0 {
		currentPath = args[0]
	}

	handlers.Open(currentPath, app)
}

func OpenFolder(args []string, app *app.Data) {
	currentPath := app.CfgFilePath
	if len(args) > 0 {
		currentPath = args[0]
	}

	handlers.Open(path.Dir(currentPath), app)
}
