package app

import (
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

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

// Commit writes a new command set and the operation that produced it as one
// unit: the snapshot first, then the history index, then the command file.
//
// The order is the contract. A failure after the snapshot leaves an orphan
// snapshot and an index that still describes the old state; a failure after the
// index leaves the history ahead of the command file, which `revert` and
// `restore` can bring back. The command file is written last because it is the
// one copy the history can reconstruct.
//
// Each file is replaced through a temporary file and a rename, so a crash
// cannot leave a half-written file behind.
func (app *Config) Commit(commands catalog.Catalog, op operation.Operation, commandName, note string) error {
	app.ReplaceCatalog(commands)
	app.History.Add(op, commandName, note)

	entry := app.History.GetLast()
	if err := app.historyFiles.WriteSnapshot(entry.File, app.AllCommandsByName); err != nil {
		return err
	}

	if err := app.SaveOnlyHistory(); err != nil {
		return err
	}

	return app.SaveCommands()
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
