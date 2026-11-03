package catalog

import "slices"

// Command is one managed command: the name it is invoked by, the aliases it
// also answers to, the line it runs, and the lines it runs instead when it is a
// Steps command.
type Command struct {
	Name    string
	Aliases []string
	Exec    string
	Steps   []string
}

// clone returns a command whose slices do not share storage with the catalog.
func (c Command) clone() Command {
	c.Aliases = slices.Clone(c.Aliases)
	c.Steps = slices.Clone(c.Steps)
	return c
}
