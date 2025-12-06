package handlers

import "github.com/fiwon123/crower/internal/data"

func AddCommand(command data.Command, app *data.App) bool {
	if app.CommandsMap.Get(command.Name) != nil {
		return false
	}

	app.CommandsMap.Add(command)
	return true
}
