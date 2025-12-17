package data

import "github.com/fiwon123/crower/pkg/crowlog"

type App struct {
	CfgFilePath        string
	LoggerInfo         crowlog.LoggerInfo
	AllCommandsByAlias CommandsMap
	AllCommandsByName  CommandsMap
}

func NewApp(cfgFilePath string, allAliases CommandsMap, allCommands CommandsMap) *App {

	return &App{
		CfgFilePath:        cfgFilePath,
		LoggerInfo:         *crowlog.New(),
		AllCommandsByAlias: allAliases,
		AllCommandsByName:  allCommands,
	}
}
