package data

type Command struct {
	Name  string
	Alias []string
	Exec  string
}

func NewCommand(name string, alias []string, exec string) *Command {
	return &Command{
		Name:  name,
		Alias: alias,
		Exec:  exec,
	}
}
