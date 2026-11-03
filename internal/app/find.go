package app

import (
	"fmt"
	"maps"
	"slices"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
)

// FindCommand resolves a key inside a history snapshot. A name wins over an
// alias, and both passes run in name order: ranging the map made a snapshot
// with a collision restore a different command on different runs.
func (app *Config) FindCommand(key string, content history.Content) (*command.Data, error) {
	allCommands, err := app.historyFiles.ReadSnapshot(content.File)
	if err != nil {
		return nil, err
	}

	names := slices.Sorted(maps.Keys(allCommands))

	for _, name := range names {
		if allCommands[name].Name == key {
			c := allCommands[name]
			return &c, nil
		}
	}

	for _, name := range names {
		if slices.Contains(allCommands[name].AllAlias, key) {
			c := allCommands[name]
			return &c, nil
		}
	}

	return nil, fmt.Errorf("Not Found")
}
