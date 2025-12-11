package data

type CommandsMap map[string]Command

func NewCommandsMap() CommandsMap {
	return make(CommandsMap)
}

func (cm CommandsMap) Get(key string) *Command {
	if key == "" {
		return nil
	}

	command, ok := cm[key]

	if !ok {
		return nil
	}
	return &command
}

func (cm CommandsMap) Add(key string, c *Command) {
	cm[key] = *c
}

func (cm CommandsMap) Remove(key string) bool {
	if cm.Get(key) == nil {
		return false
	}

	delete(cm, key)

	return true
}

func (cm CommandsMap) Update(key string, c *Command) bool {

	command := cm.Get(key)
	if command == nil {
		return false
	}

	command.Name = c.Name
	command.AllAlias = c.AllAlias
	command.Exec = c.Exec
	cm[key] = *c

	return true
}
