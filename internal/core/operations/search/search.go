package searchoperations

import (
	appdata "github.com/fiwon123/crower/internal/data/app"
	searchnotesdata "github.com/fiwon123/crower/internal/data/notes/search"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	searchhandlers "github.com/fiwon123/crower/internal/handlers/search"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
)

func SearchBrowser(args []string, app *appdata.Data) {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	searchhandlers.SearchBrowser(content, app)

	app.History.Add(operationsdata.Revert, searchnotesdata.NewSearchBrowserNote(args))
	historyhelper.Save(app)
}

func SearchFile(args []string, app *appdata.Data) {
	content, currentPath := assertArgs(args)

	out, err := searchhandlers.SearchFile(currentPath, content, app)
	assertSearchResult(out, err, app)

	app.History.Add(operationsdata.Revert, searchnotesdata.NewSearchFileNote(args))
	historyhelper.Save(app)
}

func SearchFolder(args []string, app *appdata.Data) {
	content, currentPath := assertArgs(args)

	out, err := searchhandlers.SearchFolder(currentPath, content, app)
	assertSearchResult(out, err, app)

	app.History.Add(operationsdata.Revert, searchnotesdata.NewSearchFolderNote(args))
	historyhelper.Save(app)
}

func SearchFileAndFolder(args []string, app *appdata.Data) {
	content, currentPath := assertArgs(args)

	out, err := searchhandlers.SearchFileAndFolder(currentPath, content, app)
	assertSearchResult(out, err, app)

	app.History.Add(operationsdata.Revert, searchnotesdata.NewSearchFileAndFolderNote(args))
	historyhelper.Save(app)
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
