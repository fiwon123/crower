package core

import "github.com/fiwon123/crower/internal/data/command"

func getAliasMap(commandMap command.MapData) command.MapData {
	aliasMap := command.NewMapData()

	for _, command := range commandMap {
		for _, alias := range command.AllAlias {
			aliasMap.Add(alias, &command)
		}
	}

	return aliasMap
}
