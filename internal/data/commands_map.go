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

func (cm CommandsMap) Add(key string, c Command) {
	cm[key] = c
}

func (cm CommandsMap) Remove(name string) bool {
	if cm.Get(name) == nil {
		return false
	}

	delete(cm, name)

	return true
}

func (cm CommandsMap) Update(c Command) bool {

	command := cm.Get(c.Name)
	if command == nil {
		return false
	}

	command.AllAlias = c.AllAlias
	command.Exec = c.Exec
	cm[c.Name] = c

	return true
}
