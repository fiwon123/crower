package openoperations

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func Open(args []string, app *app.Data) {

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

	handlers.Open(paths, app)

	app.History.Add(state.Open, notes.GenerateOpenNote(args))
	history.Save(app)
}

func OpenFile(args []string, app *app.Data) {
	currentPath := app.CfgFilePath
	if len(args) == 0 {
		args = append(args, currentPath)
	}

	handlers.Open(args, app)

	app.History.Add(state.Open, notes.GenerateOpenFolderNote(args))
	history.Save(app)
}

func OpenFolder(args []string, app *app.Data) {
	currentPath := app.CfgFilePath
	if len(args) == 0 {
		args = append(args, filepath.Dir(currentPath))
	}

	handlers.Open(args, app)

	app.History.Add(state.Open, notes.GenerateOpenFolderNote(args))
	history.Save(app)
}

func OpenSystem(app *app.Data) {
	err := handlers.OpenSystem(app)
	if err != nil {
		app.Logger.Error("failed to open system variable: ", "error", err)
		return
	}

	app.History.Add(state.Open, notes.GenerateOpenSystemNote())
	history.Save(app)
}
