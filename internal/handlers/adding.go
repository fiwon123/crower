package handlers

import "github.com/fiwon123/crower/internal/data"

// Add command from the cfg file.
func AddCommand(command *data.Command, app *data.App) bool {
	if app.AllCommandsByName.Get(command.Name) != nil {
		return false
	}

	for _, alias := range command.AllAlias {
		if app.AllCommandsByAlias.Get(alias) != nil || app.AllCommandsByName.Get(alias) != nil {
			return false
		}
	}

	app.AllCommandsByName.Add(command.Name, command)

	for _, alias := range command.AllAlias {
		app.AllCommandsByAlias.Add(alias, command)
	}

	return true
}
