package app

import commanddata "github.com/fiwon123/crower/internal/data/command"

func getAliasMap(commandMap commanddata.MapData) commanddata.MapData {
	aliasMap := commanddata.NewMapData()

	for _, command := range commandMap {
		for _, alias := range command.AllAlias {
			aliasMap.Add(alias, &command)
		}
	}

	return aliasMap
}
