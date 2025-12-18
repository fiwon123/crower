package handlers

import (
	"github.com/fiwon123/crower/internal/data"
)

// Update command based on the key value.
// Old values will be used if not specified in the data.Command structure.
func UpdateCommand(key string, newName string, newAlias []string, newExec string, app *data.App) bool {

	newCommand := &data.Command{
		Name:     newName,
		AllAlias: newAlias,
		Exec:     newExec,
	}

	var command *data.Command
	command = app.AllCommandsByName.Get(key)
	if command != nil {
		return performUpdate(command, newCommand, app)
	}

	command = app.AllCommandsByAlias.Get(key)
	if command != nil {
		return performUpdate(command, newCommand, app)
	}

	return false
}

func performUpdate(oldCommand *data.Command, newCommand *data.Command, app *data.App) bool {
	if newCommand.Name == "" {
		newCommand.Name = oldCommand.Name
	}

	if len(newCommand.AllAlias) == 0 {
		newCommand.AllAlias = oldCommand.AllAlias
	}

	if newCommand.Exec == "" {
		newCommand.Exec = oldCommand.Exec
	}

	if !canUpdate(newCommand, app) {
		return false
	}

	app.AllCommandsByName.Remove(oldCommand.Name)
	app.AllCommandsByName.Add(newCommand.Name, newCommand)

	for _, alias := range oldCommand.AllAlias {
		app.AllCommandsByAlias.Remove(alias)
	}

	for _, alias := range newCommand.AllAlias {
		app.AllCommandsByAlias.Add(alias, newCommand)
	}

	return true
}

func canUpdate(newCommand *data.Command, app *data.App) bool {
	if app.AllCommandsByName.Get(newCommand.Name) == nil {
		return false
	}

	return true
}
