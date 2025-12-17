package data

type Command struct {
	Name     string
	AllAlias []string
	Exec     string
}

// Create new command.
func NewCommand(name string, aliases []string, exec string) *Command {
	return &Command{
		Name:     name,
		AllAlias: aliases,
		Exec:     exec,
	}
}
