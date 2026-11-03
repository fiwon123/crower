package app

import (
	"path/filepath"

	"github.com/fiwon123/crower/pkg/utils"
)

// Save writes the current command map as a history snapshot and the history
// index next to the config file.
//
// The error is returned rather than logged: a run whose save failed did not
// persist anything, and the caller has to be able to say so with a non-zero
// exit status.
func (app *Config) Save() error {
	data := app.History.GetLast()
	if data == nil {
		return app.SaveOnlyHistory()
	}

	newDataPath := filepath.Join(app.HistoryFolderPath, data.File)
	if err := utils.CreateFileIfNotExists(newDataPath); err != nil {
		return err
	}

	if err := utils.WriteToml(app.AllCommandsByName, newDataPath); err != nil {
		return err
	}

	return utils.WriteJson(app.History, app.HistoryFilePath)
}

// SaveOnlyHistory writes the history index without taking a snapshot.
func (app *Config) SaveOnlyHistory() error {
	return utils.WriteJson(app.History, app.HistoryFilePath)
}
