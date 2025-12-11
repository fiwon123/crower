package core

import "github.com/fiwon123/crower/internal/data"

func getAliasMap(commandMap data.CommandsMap) data.CommandsMap {
	aliasMap := data.NewCommandsMap()

	for _, command := range commandMap {
		for _, alias := range command.AllAlias {
			aliasMap.Add(alias, &command)
		}
	}

	return aliasMap
}
