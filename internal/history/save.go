package history

import (
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/pkg/utils"
)

func Save(app *app.Data) {
	data := app.History.GetLast()

	newDataPath := filepath.Join(app.HistoryFolderPath, data.File)
	utils.CreateFileIfNotExists(newDataPath)
	utils.WriteToml(app.AllCommandsByName, newDataPath)

	utils.WriteJson(app.History, app.HistoryFilePath)
}
