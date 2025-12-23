package data

import (
	"fmt"
	"path/filepath"

	"github.com/fiwon123/crower/pkg/crowlog"
	"github.com/fiwon123/crower/pkg/utils"
)

type App struct {
	CfgFilePath        string
	HistoryFilePath    string
	HistoryFolderPath  string
	History            History
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

	historyFolderPath := filepath.Join(folderPath, "history")
	utils.CreateFolderIfNotExists(historyFolderPath)

	var history History
	err := utils.ReadJson(historyFilePath, &history)
	if err != nil {
		fmt.Printf("history error: %v \n", err)
	}

	return &App{
		CfgFilePath:        cfgFilePath,
		History:            history,
		HistoryFilePath:    historyFilePath,
		HistoryFolderPath:  historyFolderPath,
		LoggerInfo:         *crowlog.New(),
		OrderKeys:          orderKeys,
		AllCommandsByAlias: allAliases,
		AllCommandsByName:  allCommands,
	}
}
