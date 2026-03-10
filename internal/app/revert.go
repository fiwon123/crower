package app

import (
	"os"
	"path/filepath"

	commanddata "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

// Revert last operation
func (app *Config) RevertTo(content *history.Content) error {

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

	app.removeUntilHistory(content)

	return nil
}

func (app *Config) removeUntilHistory(content *history.Content) {

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
