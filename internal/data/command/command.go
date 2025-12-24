package command

type Data struct {
	Name     string
	AllAlias []string
	Exec     string
}

// Create new command.
func New(name string, aliases []string, exec string) *Data {

	if aliases == nil {
		aliases = []string{}
	}

	return &Data{
		Name:     name,
		AllAlias: aliases,
		Exec:     exec,
	}
}
