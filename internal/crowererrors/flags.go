package crowererrors

import (
	appdata "github.com/fiwon123/crower/internal/data/app"
)

func PrintFileAndFolderFlagsNotUsed(app *appdata.Data) {
	app.Logger.Error("file and folder flag not used")
}

func PrintNotFileAndOutputPath(app *appdata.Data) {
	app.Logger.Error("needs to specify file path and out folder")
}

func PrintNotArgs(msg string, app *appdata.Data) {
	if msg == "" {
		app.Logger.Error("need to pass arguments")
	} else {
		app.Logger.Error("need to pass arguments: ", msg)
	}
}

func PrintEmptyPaths(app *appdata.Data) {
	app.Logger.Error("empty paths")
}
