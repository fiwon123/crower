package app

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

// Initialize app based on the cfg file path.
func InitApp(cfgFilePath string) *Config {
	var orderKeys []string
	allCommands := command.NewMapData()
	var allAliases command.MapData

	if cfgFilePath != "" {
		crowerutils.CreateFolderIfNotExists(filepath.Dir(cfgFilePath))
		crowerutils.CreateFileIfNotExists(cfgFilePath)

		var err error
		orderKeys, err = crowerutils.ReadKeysTomlInOrder(cfgFilePath)
		err = crowerutils.ReadToml(cfgFilePath, &allCommands)
		if err != nil {
			panic(err)
		}
		allAliases = getAliasMap(allCommands)
	} else {
		allCommands = command.NewMapData()
		allAliases = command.NewMapData()
	}

	return NewApp(cfgFilePath, orderKeys, allAliases, allCommands)
}
