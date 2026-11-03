package app

import (
	"fmt"
	"maps"
	"slices"

	"github.com/fiwon123/crower/internal/data/command"
)

// getAliasMap builds the alias index in a stable order: commands sorted by
// name, first claim wins.
//
// Ranging the command map made a duplicate alias resolve differently between
// runs, so the same config could execute a different command each time.
func getAliasMap(commandMap command.MapData) command.MapData {
	aliasMap := command.NewMapData()

	for _, name := range slices.Sorted(maps.Keys(commandMap)) {
		c := commandMap[name]

		for _, alias := range c.AllAlias {
			if alias == "" || aliasMap.Exists(alias) {
				continue
			}

			aliasMap.Add(alias, &c)
		}
	}

	return aliasMap
}

// commandCollisions reports every key two commands claim, and every alias that
// is also another command's name. The load keeps the first claim (sorted by
// name) and lookups try the name map first; the warning names both sides so
// the user can fix the file instead of wondering why one command is unreachable.
func commandCollisions(commandMap command.MapData) []string {
	owner := map[string]string{}
	seen := map[string]bool{}
	var warnings []string

	report := func(warning string) {
		if seen[warning] {
			return
		}

		seen[warning] = true
		warnings = append(warnings, warning)
	}

	for _, name := range slices.Sorted(maps.Keys(commandMap)) {
		c := commandMap[name]

		for _, alias := range c.AllAlias {
			if alias == "" {
				continue
			}

			if existing, ok := owner[alias]; ok && existing != name {
				// A command name wins over every alias, so "using X" would be
				// false when the key is also a name; the second pass explains
				// that case.
				if _, isName := commandMap[alias]; isName {
					report(fmt.Sprintf("alias %q is claimed by %q and %q", alias, existing, name))
				} else {
					report(fmt.Sprintf("alias %q is claimed by %q and %q; using %q", alias, existing, name, existing))
				}

				continue
			}

			owner[alias] = name
		}
	}

	for _, name := range slices.Sorted(maps.Keys(commandMap)) {
		c := commandMap[name]

		for _, alias := range c.AllAlias {
			if alias == "" || alias == name {
				continue
			}

			if _, ok := commandMap[alias]; ok {
				report(fmt.Sprintf(
					"alias %q of %q is also a command name; the name wins", alias, name))
			}
		}
	}

	return warnings
}
