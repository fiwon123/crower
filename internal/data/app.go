package data

import (
	"path/filepath"

	"github.com/fiwon123/crower/pkg/crowlog"
	"github.com/fiwon123/crower/pkg/utils"
)

type App struct {
	CfgFilePath        string
	HistoryFilePath    string
	LoggerInfo         crowlog.LoggerInfo
	OrderKeys          []string
	AllCommandsByAlias CommandsMap
	AllCommandsByName  CommandsMap
}

// Create a new App containing core structures to perform all crower operations.
func NewApp(cfgFilePath string, orderKeys []string, allAliases CommandsMap, allCommands CommandsMap) *App {

	folderPath := filepath.Dir(cfgFilePath)

	historyFilePath := filepath.Join(folderPath, "history.json")
	utils.CreateFileIfNotExists(historyFilePath)

	return &App{
		CfgFilePath:        cfgFilePath,
		HistoryFilePath:    historyFilePath,
		LoggerInfo:         *crowlog.New(),
		OrderKeys:          orderKeys,
		AllCommandsByAlias: allAliases,
		AllCommandsByName:  allCommands,
	}
}
