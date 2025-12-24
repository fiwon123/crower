package commands

type MapData map[string]Data

func NewMapData() MapData {
	return make(MapData)
}

func (cm MapData) Get(key string) *Data {
	if key == "" {
		return nil
	}

	command, ok := cm[key]

	if !ok {
		return nil
	}
	return &command
}

func (cm MapData) Add(key string, c *Data) {
	cm[key] = *c
}

func (cm MapData) Remove(key string) bool {
	if cm.Get(key) == nil {
		return false
	}

	delete(cm, key)

	return true
}

func (cm MapData) Update(key string, c *Data) bool {

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
