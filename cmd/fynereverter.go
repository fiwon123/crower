package cmd

import (
	"fmt"
	"strconv"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
	"github.com/fiwon123/crower/internal/session"
)

// fyneReverter restores a history entry from the Fyne window through the same
// Session.RevertTo path the CLI's `revert` uses.
//
// The window asks the confirmation itself, so this adapter is the state change
// and not the use case with its prompt. It has no Fyne import, so the adapter
// is covered by the normal cmd tests.
type fyneReverter struct {
	cfg *session.Session
}

func newFyneReverter(cfg *session.Session) fyneReverter {
	return fyneReverter{cfg: cfg}
}

// Revert restores the command set steps entries back and records the revert.
func (r fyneReverter) Revert(steps int) (catalog.Catalog, timeline.Log, error) {
	if steps < 1 {
		return catalog.Catalog{}, timeline.Log{}, fmt.Errorf("revert needs at least one step, got %d", steps)
	}

	entry, err := r.cfg.History.GetBeforeLast(steps)
	if err != nil {
		return catalog.Catalog{}, timeline.Log{}, err
	}

	if err := r.cfg.RevertTo(entry); err != nil {
		return catalog.Catalog{}, timeline.Log{}, fmt.Errorf("revert history: %w", err)
	}

	r.cfg.History.Add(operation.Revert, "", notes.Revert([]string{strconv.Itoa(steps)}))

	if err := r.cfg.Save(); err != nil {
		return catalog.Catalog{}, timeline.Log{}, err
	}

	return r.cfg.Catalog(), r.cfg.History, nil
}
