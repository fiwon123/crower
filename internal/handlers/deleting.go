package handlers

import "github.com/fiwon123/crower/internal/data"

func DeleteCommand(name string, app *data.App) bool {
	var command *data.Command

	command = app.AllCommandsByName.Get(name)
	if command == nil {
		command = app.AllCommandsByAlias.Get(name)
		if command == nil {
			return false
		}
	}

	app.AllCommandsByName.Remove(name)

	for _, alias := range command.AllAlias {
		app.AllCommandsByAlias.Remove(alias)
	}

	return true
}
