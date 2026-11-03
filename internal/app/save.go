package app

// Save writes the current command map as a history snapshot and the history
// index next to the config file.
//
// The error is returned rather than logged: a run whose save failed did not
// persist anything, and the caller has to be able to say so with a non-zero
// exit status.
func (app *Config) Save() error {
	entry := app.History.GetLast()
	if entry == nil {
		return app.SaveOnlyHistory()
	}

	if err := app.historyFiles.WriteSnapshot(entry.File, app.AllCommandsByName); err != nil {
		return err
	}

	return app.historyFiles.SaveIndex(app.History)
}

// SaveCommands writes the command map to the configuration file.
//
// Together with RevertTo (which writes the map it restores) this is the only
// writer of the configuration file; both live in this package, so no domain has
// to know that the file is TOML or where it is.
func (app *Config) SaveCommands() error {
	return app.commands.Save(app.AllCommandsByName)
}

// DeleteHistoryBackup removes the snapshot a history entry points at.
//
// file is a name from a history entry, not a path: the storage layer checks
// that it cannot escape the history folder, because a history file can be
// hand-edited and is not trusted.
func (app *Config) DeleteHistoryBackup(file string) error {
	return app.historyFiles.DeleteSnapshot(file)
}

// SaveOnlyHistory writes the history index without taking a snapshot.
func (app *Config) SaveOnlyHistory() error {
	return app.historyFiles.SaveIndex(app.History)
}
