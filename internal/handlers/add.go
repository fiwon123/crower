package handlers

import "github.com/fiwon123/crower/internal/data"

func AddCommand(command data.Command, app *data.App) {
	app.CommandsMap.Add(&command)
}
