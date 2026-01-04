package operations

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func Open(args []string, app *app.Data) {
	currentPath := app.CfgFilePath

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
		paths = append(paths, currentPath)
	}

	handlers.Open(paths, app)
}

func OpenFolder(args []string, app *app.Data) {
	currentPath := app.CfgFilePath
	if len(args) == 0 {
		args = append(args, path.Dir(currentPath))
	}

	handlers.Open(args, app)
}

func OpenSystem(app *app.Data) {
	err := handlers.OpenSystem(app)
	if err != nil {
		fmt.Println("failed to open system variable: ", err)
		return
	}
}
