package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func SearchBrowser(args []string, app *app.Data) {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	handlers.SearchBrowser(content, app)
}

func SearchFile(args []string, app *app.Data) {
	content, currentPath := assertArgs(args)

	out, err := handlers.SearchFile(currentPath, content, app)
	assertSearchResult(out, err, app)
}

func SearchFolder(args []string, app *app.Data) {
	content, currentPath := assertArgs(args)

	out, err := handlers.SearchFolder(currentPath, content, app)
	assertSearchResult(out, err, app)
}

func SearchFileAndFolder(args []string, app *app.Data) {
	content, currentPath := assertArgs(args)

	out, err := handlers.SearchFileAndFolder(currentPath, content, app)
	assertSearchResult(out, err, app)
}

func assertArgs(args []string) (string, string) {
	content := ""
	currentPath := "."

	if len(args) > 1 {
		currentPath = args[0]
		content = args[1]
	}

	return content, currentPath
}

func assertSearchResult(out string, err error, app *app.Data) {
	if err != nil {
		app.LoggerInfo.Error("Error trying to search: ", out, err)
	}
	fmt.Println(out)
}
