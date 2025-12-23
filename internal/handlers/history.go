package handlers

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/pkg/utils"
)

func History(app *data.App) error {

	app.History.List()

	return nil
}

func SaveHistory(app *data.App) {
	data := app.History.GetLast()

	newDataPath := filepath.Join(app.HistoryFolderPath, data.File)
	utils.CreateFileIfNotExists(newDataPath)
	utils.WriteToml(app.AllCommandsByName, newDataPath)

	utils.WriteJson(app.History, app.HistoryFilePath)
}
