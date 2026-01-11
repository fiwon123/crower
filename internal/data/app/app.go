package app

import (
	"fmt"
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/pkg/crowlog"
	"github.com/fiwon123/crower/pkg/utils"
)

type Data struct {
	CfgFilePath        string
	HistoryFilePath    string
	HistoryFolderPath  string
	History            history.Data
	Logger             crowlog.Logger
	OrderKeys          []string
	AllCommandsByAlias command.MapData
	AllCommandsByName  command.MapData
}

// Create a new App containing core structures to perform all crower operations.
func New(cfgFilePath string, orderKeys []string, allAliases command.MapData, allCommands command.MapData) *Data {

	folderPath := filepath.Dir(cfgFilePath)

	historyFilePath := filepath.Join(folderPath, "history.json")
	utils.CreateFileIfNotExists(historyFilePath)

	historyFolderPath := filepath.Join(folderPath, "history")
	utils.CreateFolderIfNotExists(historyFolderPath)

	logPath := filepath.Join(folderPath, "crower.log")

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
		Logger:             *crowlog.New(logPath),
		OrderKeys:          orderKeys,
		AllCommandsByAlias: allAliases,
		AllCommandsByName:  allCommands,
	}
}
