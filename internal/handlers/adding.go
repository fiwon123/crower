package handlers

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data"
)

// Add command from the cfg file.
func AddCommand(command *data.Command, app *data.App) error {
	if command.Name == "" {
		return fmt.Errorf("empty name.")
	}

	if app.AllCommandsByName.Get(command.Name) != nil {
		return fmt.Errorf("found name, command already added.")
	}

	for _, alias := range command.AllAlias {
		if app.AllCommandsByAlias.Get(alias) != nil || app.AllCommandsByName.Get(alias) != nil {
			return fmt.Errorf("found alias, command already added.")
		}
	}

	app.AllCommandsByName.Add(command.Name, command)

	for _, alias := range command.AllAlias {
		app.AllCommandsByAlias.Add(alias, command)
	}

	return nil
}
