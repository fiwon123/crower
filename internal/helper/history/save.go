package historyhelper

import (
	"path/filepath"

	appdata "github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

// Save history data on cfg file
func Save(app *appdata.Data) {
	data := app.History.GetLast()

	newDataPath := filepath.Join(app.HistoryFolderPath, data.File)
	crowerutils.CreateFileIfNotExists(newDataPath)
	crowerutils.WriteToml(app.AllCommandsByName, newDataPath)

	crowerutils.WriteJson(app.History, app.HistoryFilePath)
}

func SaveOnlyHistory(app *appdata.Data) {
	crowerutils.WriteJson(app.History, app.HistoryFilePath)
}
