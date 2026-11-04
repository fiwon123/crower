package catalog

import "fmt"

// Add appends a command to the catalog after checking that nothing else
// already answers to its name or to one of its aliases.
//
// The rules are the ones `create` has always enforced:
//
//   - a command needs a name and an exec;
//   - the name may be neither another command's name nor its alias. A name
//     that is already an alias would shadow the command it pointed at, because
//     lookups try the name index first;
//   - an alias may be neither another command's name nor another command's
//     alias.
//
// The new command goes last, which is where a command the loader has never
// seen the section key for lands after a save and a reload.
func (c Catalog) Add(command Command) (Catalog, error) {
	if command.Name == "" {
		return c, fmt.Errorf("empty name")
	}

	if command.Exec == "" {
		return c, fmt.Errorf("empty exec")
	}

	if _, taken := c.byName[command.Name]; taken {
		return c, fmt.Errorf("found name, command already added")
	}

	if owner, taken := c.alias[command.Name]; taken {
		return c, fmt.Errorf("%q is already an alias of %q", command.Name, owner)
	}

	for _, alias := range command.Aliases {
		if alias == "" {
			// An empty alias names nothing, so it cannot shadow anything. The
			// index ignores it too.
			continue
		}

		if _, taken := c.alias[alias]; taken {
			return c, fmt.Errorf("found alias, command already added")
		}

		if _, taken := c.byName[alias]; taken {
			return c, fmt.Errorf("found alias, command already added")
		}
	}

	commands := append(c.Commands(), command)
	order := append(c.Names(), command.Name)

	return New(order, commands), nil
}
