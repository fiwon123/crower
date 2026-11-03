package catalog

import (
	"maps"
	"slices"
)

// Warnings reports every key two commands claim, and every alias that is also
// another command's name.
//
// The load keeps the first claim and lookups try the name map first, so the
// message names both sides and tells the user which one is reachable: a
// command the user cannot invoke is otherwise indistinguishable from one that
// simply did not run.
//
// The output is deduplicated and keeps the order the problems are found in, so
// a config produces the same warnings on every start.
func (c Catalog) Warnings() []string {
	var (
		warnings []string
		seen     map[string]bool
		owner    = map[string]string{}
	)

	report := func(warning string) {
		if seen == nil {
			seen = map[string]bool{}
		}

		if seen[warning] {
			return
		}

		seen[warning] = true
		warnings = append(warnings, warning)
	}

	names := slices.Sorted(maps.Keys(c.byName))

	for _, name := range names {
		for _, alias := range c.byName[name].Aliases {
			if alias == "" {
				continue
			}

			existing, claimed := owner[alias]
			if claimed && existing != name {
				// A command name wins over every alias, so saying "using X"
				// would be false when the key is also a name; the second pass
				// explains that case.
				if _, isName := c.byName[alias]; isName {
					report("alias " + quote(alias) + " is claimed by " + quote(existing) + " and " + quote(name))
				} else {
					report("alias " + quote(alias) + " is claimed by " + quote(existing) + " and " + quote(name) + "; using " + quote(existing))
				}

				continue
			}

			owner[alias] = name
		}
	}

	for _, name := range names {
		for _, alias := range c.byName[name].Aliases {
			if alias == "" || alias == name {
				continue
			}

			if _, isName := c.byName[alias]; isName {
				report("alias " + quote(alias) + " of " + quote(name) + " is also a command name; the name wins")
			}
		}
	}

	return warnings
}

// quote renders a name the way the warnings quote it.
func quote(value string) string {
	return `"` + value + `"`
}
