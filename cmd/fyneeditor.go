package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
	"github.com/fiwon123/crower/internal/session"
)

// fyneEditor applies an edit from the Fyne window through the same catalog
// rules and commit path the CLI uses. It has no Fyne import, so the adapter is
// covered by the normal cmd tests.
type fyneEditor struct {
	cfg *session.Session
}

func newFyneEditor(cfg *session.Session) fyneEditor {
	return fyneEditor{cfg: cfg}
}

// Update applies the change and returns the refreshed set and history.
func (e fyneEditor) Update(key string, changes catalog.Changes) (catalog.Catalog, timeline.Log, error) {
	result, err := e.cfg.Catalog().Update(key, changes)
	if err != nil {
		return catalog.Catalog{}, timeline.Log{}, err
	}

	note := notes.UpdateNote(result.Old, result.New, "")

	if err := e.cfg.Commit(result.Catalog, operation.Update, result.New.Name, note); err != nil {
		return catalog.Catalog{}, timeline.Log{}, err
	}

	return e.cfg.Catalog(), e.cfg.History, nil
}
