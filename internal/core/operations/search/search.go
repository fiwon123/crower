package searchoperations

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func SearchBrowser(args []string, app *app.Data) {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	handlers.SearchBrowser(content, app)

	app.History.Add(state.Revert, notes.GenerateSearchBrowserNote(args))
	history.Save(app)
}

func SearchFile(args []string, app *app.Data) {
	content, currentPath := assertArgs(args)

	out, err := handlers.SearchFile(currentPath, content, app)
	assertSearchResult(out, err, app)

	app.History.Add(state.Revert, notes.GenerateSearchFileNote(args))
	history.Save(app)
}

func SearchFolder(args []string, app *app.Data) {
	content, currentPath := assertArgs(args)

	out, err := handlers.SearchFolder(currentPath, content, app)
	assertSearchResult(out, err, app)

	app.History.Add(state.Revert, notes.GenerateSearchFolderNote(args))
	history.Save(app)
}

func SearchFileAndFolder(args []string, app *app.Data) {
	content, currentPath := assertArgs(args)

	out, err := handlers.SearchFileAndFolder(currentPath, content, app)
	assertSearchResult(out, err, app)

	app.History.Add(state.Revert, notes.GenerateSearchFileAndFolderNote(args))
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

func assertSearchResult(out string, err error, app *app.Data) {
	if err != nil {
		app.Logger.Error("Error trying to search: ", out, err)
	}

	app.Logger.Info(out)
}
