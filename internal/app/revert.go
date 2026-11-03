package app

import (
	"github.com/fiwon123/crower/internal/data/history"
)

// RevertTo restores the command map an entry points at, then drops that entry
// and every entry after it.
func (app *Config) RevertTo(entry *history.Content) error {
	commands, err := app.historyFiles.ReadSnapshot(entry.File)
	if err != nil {
		return err
	}

	if err := app.commands.Save(commands); err != nil {
		return err
	}

	app.removeUntilHistory(entry)

	return nil
}

// removeUntilHistory drops entries from the newest down to the one that was
// restored, and the snapshot file that belongs to each.
func (app *Config) removeUntilHistory(entry *history.Content) {
	for {
		last := app.History.GetLast()
		if last == nil || last.Version == entry.Version {
			return
		}

		if app.historyFiles.HasSnapshot(last.File) {
			_ = app.historyFiles.DeleteSnapshot(last.File)
		}

		app.History.RemoveLast()
	}
}
