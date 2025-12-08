package handlers

import "github.com/fiwon123/crower/internal/data"

func DeleteCommand(name string, app *data.App) bool {
	var command *data.Command

	command = app.AllCommands.Get(name)
	if command == nil {
		command = app.AllAliases.Get(name)
		if command == nil {
			return false
		}
	}

	app.AllCommands.Remove(name)

	for _, alias := range command.AllAlias {
		app.AllAliases.Remove(alias)
	}

	return true
}
