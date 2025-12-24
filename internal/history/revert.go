package history

import (
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/pkg/utils"
)

func RevertTo(content *history.Content, app *app.Data) error {

	allCommands := command.NewMapData()
	newDataPath := filepath.Join(app.HistoryFolderPath, content.File)

	err := utils.ReadToml(newDataPath, &allCommands)
	if err != nil {
		return err
	}

	err = utils.WriteToml(allCommands, app.CfgFilePath)
	if err != nil {
		return err
	}

	removeUntilHistory(content, app)

	return nil
}

func removeUntilHistory(content *history.Content, app *app.Data) {

	lastHistory := app.History.GetLast()
	for lastHistory.Version != content.Version {
		path := filepath.Join(app.HistoryFolderPath, lastHistory.File)
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
		}

		app.History.RemoveLast()
		lastHistory = app.History.GetLast()
	}

}
