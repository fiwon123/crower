package cmd

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage"
)

// The window's group icons are the file beside the configuration, so a reload
// of the same folder sees what was written.
func TestFyneGroupsSetAndReload(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))
	groups := newFyneGroups(cfg)

	if _, err := groups.SetIcon("ship", "🚀"); err != nil {
		t.Fatalf("SetIcon: %v", err)
	}

	icons, err := newFyneGroups(cfg).Icons()
	if err != nil {
		t.Fatalf("Icons: %v", err)
	}

	if icons["ship"] != "🚀" {
		t.Errorf("icons = %v, want ship", icons)
	}

	if _, err := groups.SetIcon("ship", ""); err != nil {
		t.Fatalf("clear: %v", err)
	}

	icons, err = groups.Icons()
	if err != nil {
		t.Fatalf("Icons: %v", err)
	}

	if _, found := icons["ship"]; found {
		t.Errorf("icons = %v, want the group removed", icons)
	}
}

// The file lives beside the configuration, not inside it.
func TestFyneGroupsWritesBesideTheConfiguration(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	if _, err := newFyneGroups(cfg).SetIcon("ship", "🚀"); err != nil {
		t.Fatalf("SetIcon: %v", err)
	}

	path := filepath.Join(cfg.ConfigFolder(), "groups.toml")

	if _, err := storage.NewGroups(cfg.ConfigFolder()).Load(); err != nil {
		t.Fatalf("the groups file was not written at %s: %v", path, err)
	}
}
