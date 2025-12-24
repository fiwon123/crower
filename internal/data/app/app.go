package app

import (
	"fmt"
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/commands"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/pkg/crowlog"
	"github.com/fiwon123/crower/pkg/utils"
)

type Data struct {
	CfgFilePath        string
	HistoryFilePath    string
	HistoryFolderPath  string
	History            history.Data
	LoggerInfo         crowlog.LoggerInfo
	OrderKeys          []string
	AllCommandsByAlias commands.MapData
	AllCommandsByName  commands.MapData
}

// Create a new App containing core structures to perform all crower operations.
func NewApp(cfgFilePath string, orderKeys []string, allAliases commands.MapData, allCommands commands.MapData) *Data {

	folderPath := filepath.Dir(cfgFilePath)

	historyFilePath := filepath.Join(folderPath, "history.json")
	utils.CreateFileIfNotExists(historyFilePath)

	historyFolderPath := filepath.Join(folderPath, "history")
	utils.CreateFolderIfNotExists(historyFolderPath)

	var history history.Data
	err := utils.ReadJson(historyFilePath, &history)
	if err != nil {
		fmt.Printf("history error: %v \n", err)
	}

	return &Data{
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
