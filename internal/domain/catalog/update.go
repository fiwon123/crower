package catalog

import (
	"errors"
	"fmt"
	"slices"
)

// ErrNotFound is returned when the key names neither a command nor an alias.
var ErrNotFound = errors.New("couldn't find command by name or alias")

// ErrNoChange is returned when an update names nothing to change: the command
// would be the same afterwards, so recording an operation for it would only add
// noise to the history.
var ErrNoChange = errors.New("already up-to-date")

// Changes are the fields an update may set. An empty field keeps what the
// command already has, which is why an update carries the fields to change
// rather than the whole command.
type Changes struct {
	Name    string
	Aliases []string
	Exec    string
}

// UpdateResult is a change that has not been stored yet.
type UpdateResult struct {
	// Catalog is the command set after the change.
	Catalog Catalog

	// Old and New are the command before and after, for the history note.
	Old Command
	New Command

	// Warnings name every alias the update could not claim because another
	// command owns it. They are returned rather than printed: a model does not
	// write to a stream, and the caller decides where a warning goes.
	Warnings []string
}

// Resolve returns the command a key names: a name first, then an alias.
//
// The order is the same everywhere a key is resolved. Looking only at names
// rejected every alias, which the help advertises; looking only at aliases
// would make a name that is also an alias unreachable.
func (c Catalog) Resolve(key string) (Command, bool) {
	if command, ok := c.byName[key]; ok {
		return command.clone(), true
	}

	if owner, ok := c.alias[key]; ok {
		if command, ok := c.byName[owner]; ok {
			return command.clone(), true
		}
	}

	return Command{}, false
}

// Update applies a change to the command a key resolves to.
//
// The rules are the ones an update has always had, now in one place:
//
//   - an omitted field keeps the current value, and an update that changes
//     nothing is refused rather than recorded;
//   - a name may not be another command's name or alias, and an alias may not
//     be another command's name or alias -- but only what the update actually
//     changes is validated, so a config that already carries a collision stays
//     editable;
//   - an alias list the caller supplied replaces the old one and drops a claim
//     this command did not own; an inherited list keeps such a claim dormant
//     instead of deleting data the user did not touch;
//   - the steps of a Steps command survive an update, because update has no
//     field for them (#217).
func (c Catalog) Update(key string, changes Changes) (UpdateResult, error) {
	old, found := c.Resolve(key)
	if !found {
		return UpdateResult{}, ErrNotFound
	}

	updated := old.clone()
	aliasesSupplied := len(changes.Aliases) > 0

	changed := 0

	if changes.Name != "" {
		updated.Name = changes.Name
		changed++
	}

	if aliasesSupplied {
		updated.Aliases = slices.Clone(changes.Aliases)
		changed++
	}

	if changes.Exec != "" {
		updated.Exec = changes.Exec
		changed++
	}

	if changed == 0 {
		return UpdateResult{}, ErrNoChange
	}

	// A Steps command runs its steps instead of its exec, and update has no
	// field for them. Starting from a copy of the old command is what keeps
	// them: rebuilding the command would replace a multi-line workflow with a
	// single line (#217).

	// Validate after the defaults are filled: the catalog must see the name and
	// aliases the update would actually leave behind, not the partial input.
	if err := c.canUpdate(old, updated); err != nil {
		return UpdateResult{}, err
	}

	warnings := make([]string, 0, 1)
	seen := map[string]bool{}

	warn := func(warning string) {
		if seen[warning] {
			return
		}

		seen[warning] = true
		warnings = append(warnings, warning)
	}

	kept := make([]string, 0, len(updated.Aliases))

	for _, alias := range updated.Aliases {
		owner, claimed := c.alias[alias]
		foreignAlias := claimed && owner != old.Name && owner != updated.Name

		_, isName := c.byName[alias]
		foreignName := isName && alias != old.Name && alias != updated.Name

		switch {
		case foreignAlias && aliasesSupplied:
			warn(fmt.Sprintf("alias %q belongs to %q; leaving it there", alias, owner))

		case foreignName && aliasesSupplied:
			warn(fmt.Sprintf("alias %q is the name of another command; leaving it there", alias))

		default:
			// A foreign claim on an inherited list stays in the record but
			// dormant: removing it would delete data the user did not touch.
			kept = append(kept, alias)
		}
	}

	updated.Aliases = kept

	return UpdateResult{
		Catalog:  c.replace(old.Name, updated),
		Old:      old,
		New:      updated.clone(),
		Warnings: warnings,
	}, nil
}

// canUpdate checks the final name and aliases against every other command.
//
// Only what the update changes is validated. A name or alias the command
// already holds cannot introduce a new shadow, and a config written before this
// check (or edited by hand) may already carry a collision; an update of an
// unrelated field must still work on it.
//
// Both directions are covered, because lookups try the name first: a name that
// is another command's alias, or an alias that is another command's name, would
// leave one of them unreachable.
func (c Catalog) canUpdate(old, updated Command) error {
	if updated.Name != old.Name {
		if _, taken := c.byName[updated.Name]; taken {
			return fmt.Errorf("command name already in use: %q", updated.Name)
		}

		// The owner may be this command itself when it is renamed to one of its
		// own aliases; that is not a conflict.
		if owner, claimed := c.alias[updated.Name]; claimed && owner != old.Name {
			return fmt.Errorf("%q is already an alias of %q", updated.Name, owner)
		}
	}

	for _, alias := range updated.Aliases {
		if slices.Contains(old.Aliases, alias) {
			continue
		}

		if owner, claimed := c.alias[alias]; claimed && owner != old.Name {
			return fmt.Errorf("alias already in use: %q", alias)
		}

		if _, isName := c.byName[alias]; isName && alias != old.Name {
			return fmt.Errorf("%q is already a command name", alias)
		}
	}

	return nil
}

// replace returns the catalog with one command swapped, keeping its position in
// the display order.
func (c Catalog) replace(oldName string, updated Command) Catalog {
	order := make([]string, 0, len(c.order))
	commands := make([]Command, 0, len(c.order))

	for _, name := range c.order {
		if name == oldName {
			order = append(order, updated.Name)
			commands = append(commands, updated.clone())
			continue
		}

		if name == updated.Name {
			// A rename onto a name already listed cannot happen: canUpdate
			// refuses it. Skipping keeps a duplicate out if it ever did.
			continue
		}

		order = append(order, name)
		commands = append(commands, c.byName[name].clone())
	}

	return New(order, commands)
}
