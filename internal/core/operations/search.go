package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func Search(args []string, app *app.Data) {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	handlers.Search(content, app)
}

func SearchFile(args []string, app *app.Data) {
	content := ""
	currentPath := "."
	if len(args) > 1 {
		currentPath = args[0]
		content = args[1]
	}

	out, err := handlers.SearchFile(currentPath, content, app)
	if err != nil {
		app.LoggerInfo.Error("Error trying to search: ", string(out), err)
	}
	fmt.Println(string(out))
}

func SearchFolder(args []string, app *app.Data) {
	content := ""
	currentPath := "."
	if len(args) > 1 {
		currentPath = args[0]
		content = args[1]
	}

	out, err := handlers.SearchFolder(currentPath, content, app)
	if err != nil {
		app.LoggerInfo.Error("Error trying to search: ", string(out), err)
	}
	fmt.Println(string(out))
}

func SearchFileAndFolder(args []string, app *app.Data) {
	content := ""
	currentPath := "."
	if len(args) > 1 {
		currentPath = args[0]
		content = args[1]
	}

	out, err := handlers.SearchFileAndFolder(currentPath, content, app)
	if err != nil {
		app.LoggerInfo.Error("Error trying to search: ", string(out), err)
	}
	fmt.Println(string(out))
}
