package core

import (
	"fmt"
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"

	"github.com/fiwon123/crower/pkg/utils"
)

// Initialize app based on the cfg file path.
func InitApp(cfgFilePath string) *app.Data {
	var orderKeys []string
	allCommands := command.NewMapData()
	var allAliases command.MapData

	if cfgFilePath != "" {
		utils.CreateFolderIfNotExists(filepath.Dir(cfgFilePath))
		utils.CreateFileIfNotExists(cfgFilePath)

		var err error
		orderKeys, err = utils.ReadKeysTomlInOrder(cfgFilePath)
		err = utils.ReadToml(cfgFilePath, &allCommands)
		if err != nil {
			fmt.Println("error to read toml: ", err)
		}
		allAliases = getAliasMap(allCommands)
	} else {
		allCommands = command.NewMapData()
		allAliases = command.NewMapData()
	}

	return app.New(cfgFilePath, orderKeys, allAliases, allCommands)
}
