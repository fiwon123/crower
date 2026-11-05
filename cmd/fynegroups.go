package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/session"
)

// fyneGroups stores the icon a group shows in the window.
//
// It has no Fyne import, so the adapter is covered by the normal cmd tests.
type fyneGroups struct {
	groups *storage.Groups
}

func newFyneGroups(cfg *session.Session) fyneGroups {
	return fyneGroups{groups: storage.NewGroups(cfg.ConfigFolder())}
}

// Icons returns every group that has an icon.
func (g fyneGroups) Icons() (map[string]string, error) {
	return g.groups.Load()
}

// SetIcon stores or removes one group's icon.
func (g fyneGroups) SetIcon(group, icon string) (map[string]string, error) {
	return g.groups.Set(group, icon)
}
