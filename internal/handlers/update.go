package handlers

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
)

// Update command based on the key value.
// Old values will be used if not specified in the data.Command structure.
func UpdateCommand(key string, newName string, newAlias []string, newExec string, app *app.Data) (*command.Data, *command.Data, error) {

	newCommand := command.New(newName, newAlias, newExec)

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

func performUpdate(oldCommand *command.Data, newCommand *command.Data, app *app.Data) error {

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

func canUpdate(newCommand *command.Data, app *app.Data) error {
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
