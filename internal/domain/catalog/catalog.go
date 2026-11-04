// Package catalog models the set of commands crower manages.
//
// It is the first package of the domain ring: it depends on the standard
// library only, performs no I/O, and never prints. The behaviours that decide
// what the catalog contains live here once instead of in every operation:
//
//   - the order the commands were saved in,
//   - which command owns an alias when two claim the same one (the first name
//     in sorted order wins, so a config resolves the same way on every run),
//   - the collisions the user has to be told about,
//   - and the rules an update has to satisfy.
//
// The load path used to build the alias index and the warnings itself, and two
// operations re-implemented the lookup order. Keeping them in one aggregate is
// what lets every operation share the same rules.
package catalog

import (
	"maps"
	"slices"
)

// Catalog is the set of commands, in display order, with the alias index
// derived from it.
//
// The zero value is an empty catalog and is usable.
type Catalog struct {
	order  []string
	byName map[string]Command
	alias  map[string]string
}

// New builds a catalog from commands and the order the source listed them in.
//
// The order argument is a hint: names it mentions come first, in that order,
// and every other command follows in sorted order. A name the order mentions
// but the catalog does not hold is ignored, and a repeated name is kept once --
// so a stale or hand-edited order cannot drop or duplicate a command.
//
// The alias index is deliberately built from the names in sorted order rather
// than from the display order: the owner of a contested alias must not depend
// on where the command sits in the file, or the same config would execute a
// different command after a reordering.
func New(order []string, commands []Command) Catalog {
	byName := make(map[string]Command, len(commands))

	for _, c := range commands {
		byName[c.Name] = c
	}

	seen := make(map[string]bool, len(byName))
	names := make([]string, 0, len(byName))

	for _, name := range order {
		if _, ok := byName[name]; !ok || seen[name] {
			continue
		}

		names = append(names, name)
		seen[name] = true
	}

	rest := make([]string, 0, len(byName))

	for name := range byName {
		if !seen[name] {
			rest = append(rest, name)
		}
	}

	slices.Sort(rest)
	names = append(names, rest...)

	return Catalog{
		order:  names,
		byName: byName,
		alias:  buildAliasIndex(byName),
	}
}

// buildAliasIndex maps each claimed alias to the command that owns it.
//
// The commands are visited in sorted name order and the first claim wins, which
// is the rule the loader has always used: ranging a Go map made a duplicate
// alias resolve differently between runs, so the same config could execute a
// different command each time.
func buildAliasIndex(byName map[string]Command) map[string]string {
	index := make(map[string]string)

	for _, name := range slices.Sorted(maps.Keys(byName)) {
		for _, alias := range byName[name].Aliases {
			if alias == "" {
				continue
			}

			if _, claimed := index[alias]; claimed {
				continue
			}

			index[alias] = name
		}
	}

	return index
}

// Len reports how many commands the catalog holds.
func (c Catalog) Len() int {
	return len(c.order)
}

// Names returns the command names in display order.
func (c Catalog) Names() []string {
	return slices.Clone(c.order)
}

// Commands returns the commands in display order. The values are copies, so a
// caller cannot change the catalog by writing through the result.
func (c Catalog) Commands() []Command {
	commands := make([]Command, 0, len(c.order))

	for _, name := range c.order {
		commands = append(commands, c.byName[name].clone())
	}

	return commands
}

// Get returns the command stored under a name.
func (c Catalog) Get(name string) (Command, bool) {
	command, ok := c.byName[name]
	if !ok {
		return Command{}, false
	}

	return command.clone(), true
}

// AliasOwner returns the command that owns an alias, or "" when the alias is
// unclaimed.
func (c Catalog) AliasOwner(alias string) (string, bool) {
	name, ok := c.alias[alias]
	return name, ok
}

// AliasIndex returns a copy of the alias index: alias to owning command name.
func (c Catalog) AliasIndex() map[string]string {
	return maps.Clone(c.alias)
}
