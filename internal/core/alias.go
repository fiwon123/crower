package core

import "github.com/fiwon123/crower/internal/data"

func GetAliasMap(commandMap data.CommandsMap) data.CommandsMap {
	aliasMap := data.NewCommandsMap()

	for _, command := range commandMap {
		for _, alias := range command.Alias {
			aliasMap.Add(data.NewCommand(alias, []string{}, command.Exec))
		}
	}

	return aliasMap
}
