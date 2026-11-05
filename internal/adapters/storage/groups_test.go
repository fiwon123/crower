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

// Set writes and, with an empty icon, removes.
func TestGroupsSetAndClear(t *testing.T) {
	groups := storage.NewGroups(t.TempDir())

	if _, err := groups.Set("ship", "🚀"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	icons, err := groups.Set("ship", "")
	if err != nil {
		t.Fatalf("clear: %v", err)
	}

	if _, found := icons["ship"]; found {
		t.Errorf("icons = %v, want the group removed", icons)
	}

	body, err := os.ReadFile(filepath.Join(filepath.Dir(groups.Path()), "groups.toml"))
	if err != nil {
		t.Fatalf("read groups file: %v", err)
	}

	if strings.Contains(string(body), "ship") {
		t.Errorf("groups file = %q, want the cleared group gone", body)
	}
}
