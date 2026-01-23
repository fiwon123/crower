package core

import (
	"path/filepath"

	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"

	"github.com/fiwon123/crower/pkg/crowerutils"
)

// Initialize app based on the cfg file path.
func InitApp(cfgFilePath string) *appdata.Data {
	var orderKeys []string
	allCommands := commanddata.NewMapData()
	var allAliases commanddata.MapData

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
		allCommands = commanddata.NewMapData()
		allAliases = commanddata.NewMapData()
	}

	return appdata.New(cfgFilePath, orderKeys, allAliases, allCommands)
}
