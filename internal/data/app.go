package data

type App struct {
	CfgFilePath string
	AllAliases  CommandsMap
	AllCommands CommandsMap
}

func NewApp(cfgFilePath string, allAliases CommandsMap, allCommands CommandsMap) *App {
	return &App{
		CfgFilePath: cfgFilePath,
		AllAliases:  allAliases,
		AllCommands: allCommands,
	}
}
