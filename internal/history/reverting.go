package history

import (
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/pkg/utils"
)

func RevertTo(historyData *data.HistoryData, app *data.App) error {

	allCommands := data.NewCommandsMap()
	newDataPath := filepath.Join(app.HistoryFolderPath, historyData.File)

	err := utils.ReadToml(newDataPath, &allCommands)
	if err != nil {
		return err
	}

	err = utils.WriteToml(allCommands, app.CfgFilePath)
	if err != nil {
		return err
	}

	removeUntilHistory(historyData, app)

	return nil
}

func removeUntilHistory(historyData *data.HistoryData, app *data.App) {

	lastHistory := app.History.GetLast()
	for lastHistory.Version != historyData.Version {
		path := filepath.Join(app.HistoryFolderPath, lastHistory.File)
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
		}

		app.History.RemoveLast()
		lastHistory = app.History.GetLast()
	}

}
