package handlers

import "github.com/fiwon123/crower/internal/data"

// Reset all user cfg file.
func Reset(app *data.App) {
	app.AllCommandsByName = data.NewCommandsMap()
}
