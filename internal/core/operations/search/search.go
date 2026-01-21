package searchoperations

import (
	appdata "github.com/fiwon123/crower/internal/data/app"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	searchhandlers "github.com/fiwon123/crower/internal/handlers/search"

	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func SearchBrowser(args []string, app *appdata.Data) {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	searchhandlers.SearchBrowser(content, app)

	app.History.Add(operationsdata.Revert, notes.GenerateSearchBrowserNote(args))
	history.Save(app)
}

func SearchFile(args []string, app *appdata.Data) {
	content, currentPath := assertArgs(args)

	out, err := searchhandlers.SearchFile(currentPath, content, app)
	assertSearchResult(out, err, app)

	app.History.Add(operationsdata.Revert, notes.GenerateSearchFileNote(args))
	history.Save(app)
}

func SearchFolder(args []string, app *appdata.Data) {
	content, currentPath := assertArgs(args)

	out, err := searchhandlers.SearchFolder(currentPath, content, app)
	assertSearchResult(out, err, app)

	app.History.Add(operationsdata.Revert, notes.GenerateSearchFolderNote(args))
	history.Save(app)
}

func SearchFileAndFolder(args []string, app *appdata.Data) {
	content, currentPath := assertArgs(args)

	out, err := searchhandlers.SearchFileAndFolder(currentPath, content, app)
	assertSearchResult(out, err, app)

	app.History.Add(operationsdata.Revert, notes.GenerateSearchFileAndFolderNote(args))
	history.Save(app)
}

func assertArgs(args []string) (string, string) {
	content := ""
	currentPath := "."

	if len(args) > 1 {
		content = args[0]
		currentPath = args[1]
	}

	return content, currentPath
}

func assertSearchResult(out string, err error, app *appdata.Data) {
	if err != nil {
		app.Logger.Error("Error trying to search: ", out, err)
	}

	app.Logger.Info(out)
}
