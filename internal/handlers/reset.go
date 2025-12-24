package handlers

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
)

// Reset all user cfg file.
func Reset(app *app.Data) {
	app.AllCommandsByName = command.NewMapData()
}
