package command

type MapData map[string]Data

// Create a new MapData
func NewMapData() MapData {
	return make(MapData)
}

// Get content if key exists
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

// Add new content on command map using key
func (cm MapData) Add(key string, c *Data) {
	cm[key] = *c
}

// Remove key and content from command map
func (cm MapData) Remove(key string) bool {
	if cm.Get(key) == nil {
		return false
	}

	delete(cm, key)

	return true
}

// Update command map using key and new command data
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
