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

	// Icon and Group are the list's presentation: a short marker and the set
	// the command belongs to. They change no behaviour, so a file without
	// them is complete.
	Icon  string `toml:"Icon,omitempty" json:"Icon,omitempty"`
	Group string `toml:"Group,omitempty" json:"Group,omitempty"`

	// Shell and Cwd are where the command runs: the shell it needs and the
	// folder it starts in. Empty keeps the platform default and the caller's
	// working directory, so every existing file keeps its behaviour.
	Shell string `toml:"Shell,omitempty" json:"Shell,omitempty"`
	Cwd   string `toml:"Cwd,omitempty" json:"Cwd,omitempty"`

	// Env is merged over the process environment; a value is a literal or a
	// "${VAR}" reference, never a secret. Timeout is the per-step limit as a
	// Go duration ("120s"). Both are additive.
	Env     map[string]string `toml:"Env,omitempty" json:"Env,omitempty"`
	Timeout string            `toml:"Timeout,omitempty" json:"Timeout,omitempty"`
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
