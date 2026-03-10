package crowererrors

import "github.com/fiwon123/crower/internal/app"

func PrintFileAndFolderFlagsNotUsed(app *app.Config) {
	app.Logger.Error("file and folder flag not used")
}

func PrintNotFileAndOutputPath(app *app.Config) {
	app.Logger.Error("needs to specify file path and out folder")
}

func PrintNotArgs(msg string, app *app.Config) {
	if msg == "" {
		app.Logger.Error("need to pass arguments")
	} else {
		app.Logger.Error("need to pass arguments: ", msg)
	}
}

func PrintEmptyPaths(app *app.Config) {
	app.Logger.Error("empty paths")
}
