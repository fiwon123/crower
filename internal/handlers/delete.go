package handlers

import "github.com/fiwon123/crower/internal/data"

func DeleteCommand(name string, app *data.App) bool {
	if app.CommandsMap.Get(name) == nil {
		return false
	}

	app.CommandsMap.Remove(name)

	return true
}
