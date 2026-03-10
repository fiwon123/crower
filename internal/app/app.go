package app

import (
	"path/filepath"

	commanddata "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"

	"github.com/fiwon123/crower/pkg/crowerlog"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

type Config struct {
	CfgFilePath        string
	HistoryFilePath    string
	HistoryFolderPath  string
	History            history.Data
	Logger             crowerlog.Logger
	OrderKeys          []string
	AllCommandsByAlias commanddata.MapData
	AllCommandsByName  commanddata.MapData
}

// Create a new App containing core structures to perform all crower operations.
func NewApp(cfgFilePath string, orderKeys []string, allAliases commanddata.MapData, allCommands commanddata.MapData) *Config {

	folderPath := filepath.Dir(cfgFilePath)

	historyFilePath := filepath.Join(folderPath, "historyhelper.json")
	crowerutils.CreateFileIfNotExists(historyFilePath)

	historyFolderPath := filepath.Join(folderPath, "history")
	crowerutils.CreateFolderIfNotExists(historyFolderPath)

	logPath := filepath.Join(folderPath, "crower.log")

	var history history.Data
	err := crowerutils.ReadJson(historyFilePath, &history)
	if err != nil {
		panic(err)
	}

	return &Config{
		CfgFilePath:        cfgFilePath,
		History:            history,
		HistoryFilePath:    historyFilePath,
		HistoryFolderPath:  historyFolderPath,
		Logger:             *crowerlog.New(logPath),
		OrderKeys:          orderKeys,
		AllCommandsByAlias: allAliases,
		AllCommandsByName:  allCommands,
	}
}
