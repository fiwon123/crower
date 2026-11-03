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

	oldCommand := findCommand(h.app, key)
	if oldCommand == nil {
		return nil, newCommand, fmt.Errorf("couldn't find command by name or alias")
	}

	return oldCommand, newCommand, h.performUpdate(oldCommand, newCommand)
}

// findCommand resolves a key the way every lookup does: name first, then
// alias. The input prompt and the handler share it, so the command shown for
// confirmation is the command that gets updated.
func findCommand(app *app.Config, key string) *command.Data {
	if c := app.AllCommandsByName.Get(key); c != nil {
		return c
	}

	return app.AllCommandsByAlias.Get(key)
}

func (h *Handler) performUpdate(oldCommand *command.Data, newCommand *command.Data) error {

	// An alias list the caller supplied replaces the old one; an omitted list
	// is inherited from the old command. The distinction decides whether a
	// legacy claim owned by another command is dropped (the user replaced the
	// aliases) or left dormant in the record (the user changed another field).
	aliasesSupplied := len(newCommand.AllAlias) > 0

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

	// Repoint the aliases this command owns and leave the ones another command
	// owns alone: a legacy duplicate may already have assigned them there. A
	// supplied alias list drops such a claim; an inherited one keeps it in the
	// record (dormant) instead of deleting data the user did not touch.
	kept := make([]string, 0, len(newCommand.AllAlias))
	warned := map[string]bool{}

	warn := func(warning string) {
		if warned[warning] {
			return
		}

		warned[warning] = true
		h.app.Logger.Warning(warning)
	}

	for _, alias := range newCommand.AllAlias {
		owner := h.app.AllCommandsByAlias.Get(alias)
		foreignAlias := owner != nil && owner.Name != oldCommand.Name && owner.Name != newCommand.Name
		foreignName := h.app.AllCommandsByName.Exists(alias) && alias != oldCommand.Name && alias != newCommand.Name

		switch {
		case foreignAlias && aliasesSupplied:
			warn(fmt.Sprintf("alias %q belongs to %q; leaving it there", alias, owner.Name))
			continue
		case foreignName && aliasesSupplied:
			warn(fmt.Sprintf("alias %q is the name of another command; leaving it there", alias))
			continue
		case foreignAlias || foreignName:
			kept = append(kept, alias)
		default:
			kept = append(kept, alias)
			h.app.AllCommandsByAlias.Add(alias, newCommand)
		}
	}

	newCommand.AllAlias = kept

	// Drop the old aliases no longer kept, but only entries this command owns.
	for _, alias := range oldCommand.AllAlias {
		if slices.Contains(kept, alias) {
			continue
		}

		if owner := h.app.AllCommandsByAlias.Get(alias); owner != nil && owner.Name == oldCommand.Name {
			h.app.AllCommandsByAlias.Remove(alias)
		}
	}

	h.app.AllCommandsByName.Remove(oldCommand.Name)
	h.app.AllCommandsByName.Add(newCommand.Name, newCommand)

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
