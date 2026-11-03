package updatedom

import (
	"fmt"
	"slices"

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

	// Validate after the defaults are filled: the map must see the name and
	// aliases the update would actually leave behind, not the partial input.
	if err := h.canUpdate(oldCommand, newCommand); err != nil {
		return err
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

// canUpdate checks the final name and aliases against every other command.
//
// Only what the update changes is validated. A name or alias the command
// already holds cannot introduce a new shadow, and a config written before this
// check (or edited by hand) may already carry a collision; an update of an
// unrelated field must still work on it.
//
// The checks cover both directions, because lookups try the name map first: a
// name that is another command's alias, or an alias that is another command's
// name, would leave one of them unreachable.
func (h *Handler) canUpdate(oldCommand *command.Data, newCommand *command.Data) error {
	if newCommand.Name != oldCommand.Name {
		if h.app.AllCommandsByName.Get(newCommand.Name) != nil {
			return fmt.Errorf("command name already in use: %q", newCommand.Name)
		}

		// The owner may be this command itself when it is renamed to one of
		// its own aliases; that is not a conflict.
		if owner := h.app.AllCommandsByAlias.Get(newCommand.Name); owner != nil && owner.Name != oldCommand.Name {
			return fmt.Errorf("%q is already an alias of %q", newCommand.Name, owner.Name)
		}
	}

	for _, alias := range newCommand.AllAlias {
		if slices.Contains(oldCommand.AllAlias, alias) {
			continue
		}

		if owner := h.app.AllCommandsByAlias.Get(alias); owner != nil && owner.Name != oldCommand.Name {
			return fmt.Errorf("alias already in use: %q", alias)
		}

		if owner := h.app.AllCommandsByName.Get(alias); owner != nil && owner.Name != oldCommand.Name {
			return fmt.Errorf("%q is already a command name", alias)
		}
	}

	return nil
}
