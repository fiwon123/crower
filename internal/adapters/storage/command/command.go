package command

type Data struct {
	Name     string
	AllAlias []string
	Exec     string
	Steps    []string `toml:"Steps,omitempty" json:"Steps,omitempty"`

	// Description says what the command is for; Params names its positional
	// values. Both are metadata: they describe the command and change no
	// behaviour, so an old crower reading them ignores what it does not know.
	Description string   `toml:"Description,omitempty" json:"Description,omitempty"`
	Params      []string `toml:"Params,omitempty" json:"Params,omitempty"`
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
