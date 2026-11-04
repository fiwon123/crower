package command

type Data struct {
	Name     string
	AllAlias []string
	Exec     string
	Steps    []string `toml:"Steps,omitempty" json:"Steps,omitempty"`
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
