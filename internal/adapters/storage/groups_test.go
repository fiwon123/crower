package storage_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage"
)

// A missing file is "no group has an icon", not an error.
func TestGroupsLoadIsEmptyOnAFreshInstallation(t *testing.T) {
	icons, err := storage.NewGroups(t.TempDir()).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(icons) != 0 {
		t.Errorf("icons = %v, want empty", icons)
	}
}

// Icons survive a save and a load, including a group name that would become a
// nested table if the file were a flat map.
func TestGroupsRoundTripIncludingADottedName(t *testing.T) {
	groups := storage.NewGroups(t.TempDir())

	if err := groups.Save(map[string]string{"ship": "🚀", "v1.2": "🏷"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	icons, err := groups.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if icons["ship"] != "🚀" || icons["v1.2"] != "🏷" {
		t.Errorf("icons = %v, want both entries", icons)
	}
}

// Set creates the group, and clearing the icon keeps it: a group is a record
// the window can offer before any command uses it.
func TestGroupsSetAndClearKeepsTheGroup(t *testing.T) {
	groups := storage.NewGroups(t.TempDir())

	if _, err := groups.Set("ship", "🚀"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	icons, err := groups.Set("ship", "")
	if err != nil {
		t.Fatalf("clear: %v", err)
	}

	if icon, found := icons["ship"]; !found || icon != "" {
		t.Errorf("icons = %v, want the group kept without an icon", icons)
	}

	body, err := os.ReadFile(filepath.Join(filepath.Dir(groups.Path()), "groups.toml"))
	if err != nil {
		t.Fatalf("read groups file: %v", err)
	}

	if !strings.Contains(string(body), "ship") {
		t.Errorf("groups file = %q, want the group kept", body)
	}
}

// A group without an icon survives a save and a load.
func TestGroupsKeepAGroupWithNoIcon(t *testing.T) {
	groups := storage.NewGroups(t.TempDir())

	if err := groups.Save(map[string]string{"deploy": "", "ship": "🚀"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	icons, err := groups.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if icon, found := icons["deploy"]; !found || icon != "" {
		t.Errorf("icons = %v, want deploy with no icon", icons)
	}
}

// Remove deletes a group from both the map and the file.
func TestGroupsRemove(t *testing.T) {
	groups := storage.NewGroups(t.TempDir())

	if _, err := groups.Set("deploy", "🚀"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	icons, err := groups.Remove("deploy")
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}

	if _, found := icons["deploy"]; found {
		t.Errorf("icons = %v, want the group removed", icons)
	}
}
