package core

import "github.com/fiwon123/crower/internal/data/commands"

func getAliasMap(commandMap commands.MapData) commands.MapData {
	aliasMap := commands.NewMapData()

	for _, command := range commandMap {
		for _, alias := range command.AllAlias {
			aliasMap.Add(alias, &command)
		}
	}

	return aliasMap
}
