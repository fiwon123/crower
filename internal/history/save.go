package history

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

// Save history data on cfg file
func Save(app *app.Data) {
	data := app.History.GetLast()

	newDataPath := filepath.Join(app.HistoryFolderPath, data.File)
	crowerutils.CreateFileIfNotExists(newDataPath)
	crowerutils.WriteToml(app.AllCommandsByName, newDataPath)

	crowerutils.WriteJson(app.History, app.HistoryFilePath)
}

func SaveOnlyHistory(app *app.Data) {
	crowerutils.WriteJson(app.History, app.HistoryFilePath)
}
