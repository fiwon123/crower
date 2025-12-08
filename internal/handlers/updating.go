package handlers

import "github.com/fiwon123/crower/internal/data"

func UpdateCommand(command data.Command, app *data.App) bool {
	return app.AllCommands.Update(command)
}
