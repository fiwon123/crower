package handlers

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
)

// Delete command from the cfg file.
func DeleteCommand(key string, app *app.Data) (*command.Data, bool) {
	command := app.AllCommandsByName.Get(key)
	if command == nil {
		command = app.AllCommandsByAlias.Get(key)
		if command == nil {
			return nil, false
		}
	}

	app.AllCommandsByName.Remove(key)

	for _, alias := range command.AllAlias {
		app.AllCommandsByAlias.Remove(alias)
	}

	return command, true
}
