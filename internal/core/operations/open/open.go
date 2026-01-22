package openoperations

import (
	"path/filepath"

	appdata "github.com/fiwon123/crower/internal/data/app"
	opennotesdata "github.com/fiwon123/crower/internal/data/notes/open"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	openhandlers "github.com/fiwon123/crower/internal/handlers/open"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
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

	app.History.Add(operationsdata.Open, opennotesdata.NewOpenNote(args))
	historyhelper.Save(app)
}

func OpenFile(args []string, app *appdata.Data) {
	currentPath := app.CfgFilePath
	if len(args) == 0 {
		args = append(args, currentPath)
	}

	openhandlers.Open(args, app)

	app.History.Add(operationsdata.Open, opennotesdata.NewOpenFolderNote(args))
	historyhelper.Save(app)
}

func OpenFolder(args []string, app *appdata.Data) {
	currentPath := app.CfgFilePath
	if len(args) == 0 {
		args = append(args, filepath.Dir(currentPath))
	}

	openhandlers.Open(args, app)

	app.History.Add(operationsdata.Open, opennotesdata.NewOpenFolderNote(args))
	historyhelper.Save(app)
}

func OpenSystem(app *appdata.Data) {
	err := openhandlers.OpenSystem(app)
	if err != nil {
		app.Logger.Error("failed to open system variable: ", "error", err)
		return
	}

	app.History.Add(operationsdata.Open, opennotesdata.NewOpenSystemNote())
	historyhelper.Save(app)
}
