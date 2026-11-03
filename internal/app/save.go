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

	newDataPath := filepath.Join(app.historyFolderPath, data.File)
	if err := utils.CreateFileIfNotExists(newDataPath); err != nil {
		return err
	}

	if err := utils.WriteToml(app.AllCommandsByName, newDataPath); err != nil {
		return err
	}

	return utils.WriteJson(app.History, app.historyFilePath)
}

// SaveCommands writes the command map to the configuration file.
//
// It is the only place that knows the configuration is a TOML file at a
// particular path. create, update, delete and reset used to write it
// themselves, so a change to how state is stored had five sites.
func (app *Config) SaveCommands() error {
	return utils.WriteToml(app.AllCommandsByName, app.cfgFilePath)
}

// DeleteHistoryBackup removes the snapshot a history entry points at.
func (app *Config) DeleteHistoryBackup(file string) error {
	return utils.DeleteFile(filepath.Join(app.historyFolderPath, file))
}

// SaveOnlyHistory writes the history index without taking a snapshot.
func (app *Config) SaveOnlyHistory() error {
	return utils.WriteJson(app.History, app.historyFilePath)
}
