package app

import (
	"fmt"
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/pkg/utils"
)

// InitApp reads the configuration at cfgFilePath and builds a run from it.
func InitApp(cfgFilePath string, opts Options) (*Config, error) {
	opts = opts.withDefaults()

	allCommands := command.NewMapData()
	allAliases := command.NewMapData()
	var orderKeys []string

	if cfgFilePath != "" {
		if err := utils.CreateFolderIfNotExists(filepath.Dir(cfgFilePath)); err != nil {
			return nil, fmt.Errorf("prepare config folder: %w", err)
		}

		if err := utils.CreateFileIfNotExists(cfgFilePath); err != nil {
			return nil, fmt.Errorf("prepare config file %s: %w", cfgFilePath, err)
		}

		// Both reads are checked. The previous code assigned this error and then
		// overwrote it before testing it, so a config whose key order could not
		// be read lost the order of every command without a word.
		keys, err := utils.ReadKeysTomlInOrder(cfgFilePath)
		if err != nil {
			return nil, fmt.Errorf("read command order from %s: %w", cfgFilePath, err)
		}
		orderKeys = keys

		if err := utils.ReadToml(cfgFilePath, &allCommands); err != nil {
			return nil, fmt.Errorf("read config %s: %w", cfgFilePath, err)
		}

		allAliases = getAliasMap(allCommands)
	}

	return NewApp(cfgFilePath, orderKeys, allAliases, allCommands, opts)
}
