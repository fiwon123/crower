package history

import (
	"os"
	"path/filepath"

	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	historydata "github.com/fiwon123/crower/internal/data/history"

	"github.com/fiwon123/crower/pkg/crowerutils"
)

// Revert last operation
func RevertTo(content *historydata.Content, app *appdata.Data) error {

	allCommands := commanddata.NewMapData()
	newDataPath := filepath.Join(app.HistoryFolderPath, content.File)

	err := crowerutils.ReadToml(newDataPath, &allCommands)
	if err != nil {
		return err
	}

	err = crowerutils.WriteToml(allCommands, app.CfgFilePath)
	if err != nil {
		return err
	}

	removeUntilHistory(content, app)

	return nil
}

func removeUntilHistory(content *historydata.Content, app *appdata.Data) {

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
