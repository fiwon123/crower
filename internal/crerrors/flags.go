package crerrors

import (
	"github.com/fiwon123/crower/internal/data/app"
)

func PrintFileAndFolderFlagsNotUsed(app *app.Data) {
	app.Logger.Error("file and folder flag not used")
}

func PrintNotFileAndOutputPath(app *app.Data) {
	app.Logger.Error("needs to specify file path and out folder")
}

func PrintNotArgs(msg string, app *app.Data) {
	if msg == "" {
		app.Logger.Error("need to pass arguments")
	} else {
		app.Logger.Error("need to pass arguments: ", msg)
	}
}

func PrintEmptyPaths(app *app.Data) {
	app.Logger.Error("empty paths")
}
