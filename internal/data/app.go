package data

type App struct {
	CfgFilePath        string
	AllCommandsByAlias CommandsMap
	AllCommandsByName  CommandsMap
}

func NewApp(cfgFilePath string, allAliases CommandsMap, allCommands CommandsMap) *App {
	return &App{
		CfgFilePath:        cfgFilePath,
		AllCommandsByAlias: allAliases,
		AllCommandsByName:  allCommands,
	}
}
