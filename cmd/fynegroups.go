package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/session"
)

// fyneGroups stores the icon a group shows in the window.
//
// It has no Fyne import, so the adapter is covered by the normal cmd tests.
type fyneGroups struct {
	cfg *session.Session
}

func newFyneGroups(cfg *session.Session) fyneGroups {
	return fyneGroups{cfg: cfg}
}

// Icons returns every stored group with its icon.
//
// The file is looked up per call rather than captured: switching a setup moves
// the configuration folder, and the icons have to move with it.
func (g fyneGroups) Icons() (map[string]string, error) {
	return storage.NewGroups(g.cfg.ConfigFolder()).Load()
}

// SetIcon stores one group and its icon, creating the group when it is new.
func (g fyneGroups) SetIcon(group, icon string) (map[string]string, error) {
	return storage.NewGroups(g.cfg.ConfigFolder()).Set(group, icon)
}

// Remove deletes one group.
func (g fyneGroups) Remove(group string) (map[string]string, error) {
	return storage.NewGroups(g.cfg.ConfigFolder()).Remove(group)
}
