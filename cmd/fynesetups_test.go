package cmd

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/catalog"
)

// The switcher lists the original installation first and the setups sorted, so
// the order does not depend on the filesystem.
func TestSetupNamesListsTheDefaultFirst(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	for _, name := range []string{"work", "alpha"} {
		if err := os.MkdirAll(filepath.Join(home, "crower", "setups", name), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", name, err)
		}
	}

	got := setupNames()
	want := []string{"default", "alpha", "work"}

	if !slices.Equal(got, want) {
		t.Errorf("setupNames() = %v, want %v", got, want)
	}
}

// Switching reloads the same session onto the other setup and writes the marker
// `crower setup use` writes, so the next CLI run opens the same one.
func TestFyneSetupsSwitchReloadsTheSessionAndPersistsTheChoice(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CROWER_SETUP", "")

	if err := rootCmd.PersistentFlags().Set(flagSetup, ""); err != nil {
		t.Fatalf("reset --setup: %v", err)
	}

	t.Cleanup(func() {
		_ = rootCmd.PersistentFlags().Set(flagSetup, "")
	})

	defaultPath := filepath.Join(home, "crower", "crower.yaml")
	cfg := loadFyneConfigAt(t, defaultPath)
	cfg.ReplaceCatalog(catalog.New(nil, []catalog.Command{{Name: "home", Exec: "echo home"}}))

	if err := cfg.SaveCommands(); err != nil {
		t.Fatalf("save the default setup: %v", err)
	}

	workPath := filepath.Join(home, "crower", "setups", "work", "crower.yaml")
	work := loadFyneConfigAt(t, workPath)
	work.ReplaceCatalog(catalog.New(nil, []catalog.Command{{Name: "office", Exec: "echo office"}}))

	if err := work.SaveCommands(); err != nil {
		t.Fatalf("save the work setup: %v", err)
	}

	setups := newFyneSetups(cfg)

	state, err := setups.Switch("work")
	if err != nil {
		t.Fatalf("Switch work: %v", err)
	}

	if _, found := state.Catalog.Get("office"); !found {
		t.Error("the switched state does not hold the work command")
	}

	if _, found := state.Catalog.Get("home"); found {
		t.Error("the switched state still holds the default command")
	}

	if cfg.ConfigFile() != workPath {
		t.Errorf("ConfigFile() = %q, want %q", cfg.ConfigFile(), workPath)
	}

	marker, err := os.ReadFile(filepath.Join(home, "crower", "setup"))
	if err != nil {
		t.Fatalf("read the setup marker: %v", err)
	}

	if strings.TrimSpace(string(marker)) != "work" {
		t.Errorf("marker = %q, want work", marker)
	}

	// "default" is the original installation, not a folder called default.
	if _, err := setups.Switch("default"); err != nil {
		t.Fatalf("Switch default: %v", err)
	}

	if cfg.ConfigFile() != defaultPath {
		t.Errorf("ConfigFile() = %q, want the default installation", cfg.ConfigFile())
	}
}
