package app

import (
	"fmt"
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
// Together with RevertTo (which writes the map it restores) this is the only
// writer of the configuration file; both live in this package, so no domain has
// to know that the file is TOML or where it is. create, update, delete and
// reset used to write it themselves, which made the storage format a
// five-site change.
func (app *Config) SaveCommands() error {
	return utils.WriteToml(app.AllCommandsByName, app.cfgFilePath)
}

// DeleteHistoryBackup removes the snapshot a history entry points at.
//
// file is a name from a history entry, not a path: it must not escape the
// history folder. A history file can be hand-edited, so the name is checked
// here rather than trusted from the caller.
func (app *Config) DeleteHistoryBackup(file string) error {
	if file == "" || filepath.Base(file) != file {
		return fmt.Errorf("invalid history file name %q", file)
	}

	return utils.DeleteFile(filepath.Join(app.historyFolderPath, file))
}

// SaveOnlyHistory writes the history index without taking a snapshot.
func (app *Config) SaveOnlyHistory() error {
	return utils.WriteJson(app.History, app.historyFilePath)
}
