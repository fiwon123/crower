package openoperations

import (
	"path/filepath"

	appdata "github.com/fiwon123/crower/internal/data/app"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	openhandlers "github.com/fiwon123/crower/internal/handlers/open"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func Open(args []string, app *appdata.Data) {

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

	openhandlers.Open(paths, app)

	app.History.Add(operationsdata.Open, notes.GenerateOpenNote(args))
	history.Save(app)
}

func OpenFile(args []string, app *appdata.Data) {
	currentPath := app.CfgFilePath
	if len(args) == 0 {
		args = append(args, currentPath)
	}

	openhandlers.Open(args, app)

	app.History.Add(operationsdata.Open, notes.GenerateOpenFolderNote(args))
	history.Save(app)
}

func OpenFolder(args []string, app *appdata.Data) {
	currentPath := app.CfgFilePath
	if len(args) == 0 {
		args = append(args, filepath.Dir(currentPath))
	}

	openhandlers.Open(args, app)

	app.History.Add(operationsdata.Open, notes.GenerateOpenFolderNote(args))
	history.Save(app)
}

func OpenSystem(app *appdata.Data) {
	err := openhandlers.OpenSystem(app)
	if err != nil {
		app.Logger.Error("failed to open system variable: ", "error", err)
		return
	}

	app.History.Add(operationsdata.Open, notes.GenerateOpenSystemNote())
	history.Save(app)
}
