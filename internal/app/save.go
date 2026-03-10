package app

import (
	"path/filepath"

	"github.com/fiwon123/crower/pkg/utils"
)

// Save history data on cfg file
func (app *Config) Save() {
	data := app.History.GetLast()

	newDataPath := filepath.Join(app.HistoryFolderPath, data.File)
	utils.CreateFileIfNotExists(newDataPath)
	utils.WriteToml(app.AllCommandsByName, newDataPath)

	utils.WriteJson(app.History, app.HistoryFilePath)
}

func (app *Config) SaveOnlyHistory() {
	utils.WriteJson(app.History, app.HistoryFilePath)
}
