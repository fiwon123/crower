package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
	"github.com/fiwon123/crower/internal/errors"
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

// Create adds a command and returns the refreshed set and history.
func (e fyneEditor) Create(name string, aliases []string, exec string) (catalog.Catalog, timeline.Log, error) {
	command := catalog.Command{Name: name, Aliases: aliases, Exec: exec}

	next, err := e.cfg.Catalog().Add(command)
	if err != nil {
		return catalog.Catalog{}, timeline.Log{}, err
	}

	if err := e.cfg.Commit(next, operation.Create, command.Name, notes.CreateCommandNote(command)); err != nil {
		return catalog.Catalog{}, timeline.Log{}, err
	}

	return e.cfg.Catalog(), e.cfg.History, nil
}

// Delete removes a command and returns the refreshed set and history.
func (e fyneEditor) Delete(key string) (catalog.Catalog, timeline.Log, error) {
	next, removed, found := e.cfg.Catalog().Remove(key)
	if !found {
		return catalog.Catalog{}, timeline.Log{}, errors.CommandNotFound(key)
	}

	note := notes.DeleteCommandNote(removed, []string{removed.Name})

	if err := e.cfg.Commit(next, operation.Delete, removed.Name, note); err != nil {
		return catalog.Catalog{}, timeline.Log{}, err
	}

	return e.cfg.Catalog(), e.cfg.History, nil
}
