package data

type CommandsMap map[string]*Command

func NewCommandsMap() CommandsMap {
	return make(CommandsMap)
}

func (cm CommandsMap) Add(c *Command) {
	cm[c.Name] = c
}

func (cm CommandsMap) Remove(name string) {
	delete(cm, name)
}

func (cm CommandsMap) Update(c *Command) {
	if command, ok := cm[c.Name]; ok {
		command.Alias = c.Alias
		command.Exec = c.Exec
	}
}
