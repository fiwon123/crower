package resethandlers

import (
	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
)

// Reset all user cfg file.
func Reset(app *appdata.Data) {
	app.AllCommandsByName = commanddata.NewMapData()
}
