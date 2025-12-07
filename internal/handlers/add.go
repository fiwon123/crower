package handlers

import "github.com/fiwon123/crower/internal/data"

func AddCommand(command data.Command, app *data.App) bool {
	if app.CommandsMap.Get(command.Name) != nil {
		return false
	}

	for _, alias := range command.Alias {
		if app.AliasMap.Get(alias) != nil || app.CommandsMap.Get(alias) != nil {
			return false
		}
	}

	app.CommandsMap.Add(command)

	for _, alias := range command.Alias {
		app.AliasMap.Add(*data.NewCommand(alias, []string{}, command.Exec))
	}

	return true
}
