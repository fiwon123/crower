package catalog

import (
	"maps"
	"slices"
)

// Command is one managed command: the name it is invoked by, the aliases it
// also answers to, the line it runs, and the lines it runs instead when it is a
// Steps command.
type Command struct {
	Name    string
	Aliases []string
	Exec    string
	Steps   []string

	// Description and Params are metadata: what the command is for and what
	// its positional values mean. They never change how it runs.
	Description string
	Params      []string

	// Shell names the shell the command runs through ("sh", "pwsh", "cmd"),
	// and Cwd is the folder it runs in. Empty Shell keeps the platform's own
	// shell; empty Cwd keeps the caller's working directory. Every step of a
	// Steps command inherits both.
	Shell string
	Cwd   string

	// Env is merged over the process environment, and Timeout is the per-step
	// limit ("120s"). Empty Env adds nothing; empty Timeout means no limit.
	// Every step inherits both.
	Env     map[string]string
	Timeout string
}

// clone returns a command whose slices do not share storage with the catalog.
func (c Command) clone() Command {
	c.Aliases = slices.Clone(c.Aliases)
	c.Steps = slices.Clone(c.Steps)
	c.Params = slices.Clone(c.Params)
	c.Env = maps.Clone(c.Env)
	return c
}
