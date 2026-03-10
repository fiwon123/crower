package updatedom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
)

type Handler struct {
	app *app.Config
}

func NewHandler(app *app.Config) *Handler {

	return &Handler{
		app: app,
	}
}

// Update command based on the key value.
// Old values will be used if not specified in the data.Command structure.
func (h *Handler) UpdateCommand(key string, newName string, newAlias []string, newExec string) (*command.Data, *command.Data, error) {

	newCommand := command.New(newName, newAlias, newExec)

	oldCommand := h.app.AllCommandsByName.Get(key)
	if oldCommand != nil {
		return oldCommand, newCommand, h.performUpdate(oldCommand, newCommand)
	}

	oldCommand = h.app.AllCommandsByAlias.Get(key)
	if oldCommand != nil {
		return oldCommand, newCommand, h.performUpdate(oldCommand, newCommand)
	}

	return oldCommand, newCommand, fmt.Errorf("couldn't find command by name or alias")
}

func (h *Handler) performUpdate(oldCommand *command.Data, newCommand *command.Data) error {

	err := h.canUpdate(newCommand)
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

	h.app.AllCommandsByName.Remove(oldCommand.Name)
	h.app.AllCommandsByName.Add(newCommand.Name, newCommand)

	for _, alias := range oldCommand.AllAlias {
		h.app.AllCommandsByAlias.Remove(alias)
	}

	for _, alias := range newCommand.AllAlias {
		h.app.AllCommandsByAlias.Add(alias, newCommand)
	}

	return nil
}

func (h *Handler) canUpdate(newCommand *command.Data) error {
	if h.app.AllCommandsByName.Get(newCommand.Name) != nil {
		return fmt.Errorf("command name already in use: %v", newCommand.Name)
	}

	for _, alias := range newCommand.AllAlias {
		if h.app.AllCommandsByAlias.Get(alias) != nil {
			return fmt.Errorf("alias already in use: %v", alias)
		}
	}

	return nil
}
