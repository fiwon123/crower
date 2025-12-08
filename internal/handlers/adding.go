package handlers

import "github.com/fiwon123/crower/internal/data"

func AddCommand(command data.Command, app *data.App) bool {
	if app.AllCommands.Get(command.Name) != nil {
		return false
	}

	for _, alias := range command.AllAlias {
		if app.AllAliases.Get(alias) != nil || app.AllCommands.Get(alias) != nil {
			return false
		}
	}

	app.AllCommands.Add(command.Name, command)

	for _, alias := range command.AllAlias {
		app.AllAliases.Add(alias, command)
	}

	return true
}
