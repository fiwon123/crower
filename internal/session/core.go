package session

import (
	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/adapters/storage/command"
)

// Load reads the configuration at cfgFilePath and builds a run from it.
func Load(cfgFilePath string, opts Options) (*Session, error) {
	opts = opts.withDefaults()

	allCommands := command.NewMapData()
	allAliases := command.NewMapData()
	var orderKeys []string

	if cfgFilePath != "" {
		loaded, order, err := storage.NewCommandFile(cfgFilePath).Load()
		if err != nil {
			return nil, err
		}

		allCommands = loaded
		orderKeys = order
		allAliases = getAliasMap(allCommands)
	}

	return New(cfgFilePath, orderKeys, allAliases, allCommands, opts)
}
