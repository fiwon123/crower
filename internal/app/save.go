package app

import (
	"path/filepath"

	"github.com/fiwon123/crower/pkg/crowerutils"
)

// Save history data on cfg file
func (app *Config) Save() {
	data := app.History.GetLast()

	newDataPath := filepath.Join(app.HistoryFolderPath, data.File)
	crowerutils.CreateFileIfNotExists(newDataPath)
	crowerutils.WriteToml(app.AllCommandsByName, newDataPath)

	crowerutils.WriteJson(app.History, app.HistoryFilePath)
}

func (app *Config) SaveOnlyHistory() {
	crowerutils.WriteJson(app.History, app.HistoryFilePath)
}
