package handlers

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data"
)

// Update command based on the key value.
// Old values will be used if not specified in the data.Command structure.
func UpdateCommand(key string, newName string, newAlias []string, newExec string, app *data.App) (*data.Command, *data.Command, error) {

	newCommand := &data.Command{
		Name:     newName,
		AllAlias: newAlias,
		Exec:     newExec,
	}

	oldCommand := app.AllCommandsByName.Get(key)
	if oldCommand != nil {
		return oldCommand, newCommand, performUpdate(oldCommand, newCommand, app)
	}

	oldCommand = app.AllCommandsByAlias.Get(key)
	if oldCommand != nil {
		return oldCommand, newCommand, performUpdate(oldCommand, newCommand, app)
	}

	return oldCommand, newCommand, fmt.Errorf("couldn't find command by name or alias")
}

func performUpdate(oldCommand *data.Command, newCommand *data.Command, app *data.App) error {

	err := canUpdate(newCommand, app)
	if err != nil {
		return err
	}

	countFieldsUpdate := 0
	if newCommand.Name == "" {
		newCommand.Name = oldCommand.Name
	} else {
		countFieldsUpdate += 1
	}

	if len(newCommand.AllAlias) == 0 {
		newCommand.AllAlias = oldCommand.AllAlias
	} else {
		countFieldsUpdate += 1
	}

	if newCommand.Exec == "" {
		newCommand.Exec = oldCommand.Exec
	} else {
		countFieldsUpdate += 1
	}

	if countFieldsUpdate == 0 {
		return fmt.Errorf("already up-to-date")
	}

	app.AllCommandsByName.Remove(oldCommand.Name)
	app.AllCommandsByName.Add(newCommand.Name, newCommand)

	for _, alias := range oldCommand.AllAlias {
		app.AllCommandsByAlias.Remove(alias)
	}

	for _, alias := range newCommand.AllAlias {
		app.AllCommandsByAlias.Add(alias, newCommand)
	}

	return nil
}

func canUpdate(newCommand *data.Command, app *data.App) error {
	if app.AllCommandsByName.Get(newCommand.Name) != nil {
		return fmt.Errorf("command name already in use: %v", newCommand.Name)
	}

	for _, alias := range newCommand.AllAlias {
		if app.AllCommandsByAlias.Get(alias) != nil {
			return fmt.Errorf("alias already in use: %v", alias)
		}
	}

	return nil
}
