package handlers

import "github.com/fiwon123/crower/internal/data"

func Reset(app *data.App) {
	app.AllCommands = data.NewCommandsMap()
}
