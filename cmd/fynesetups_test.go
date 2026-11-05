package cmd

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/catalog"
)

// A setup's command count comes from its own configuration file, and looking at
// a setup that does not exist creates nothing.
func TestFyneSetupsCountsTheCommandsInASetup(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	workPath := filepath.Join(home, "crower", "setups", "work", "crower.yaml")
	work := loadFyneConfigAt(t, workPath)
	work.ReplaceCatalog(catalog.New(nil, []catalog.Command{
		{Name: "deploy", Exec: "make deploy"},
		{Name: "build", Exec: "make build"},
	}))

	if err := work.SaveCommands(); err != nil {
		t.Fatalf("save work: %v", err)
	}

	setups := fyneSetups{current: "default"}

	count, err := setups.Commands("work")
	if err != nil {
		t.Fatalf("Commands work: %v", err)
	}

	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}

	count, err = setups.Commands("missing")
	if err != nil {
		t.Fatalf("Commands missing: %v", err)
	}

	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}

	if _, err := os.Stat(filepath.Join(home, "crower", "setups", "missing")); !os.IsNotExist(err) {
		t.Errorf("counting a missing setup created it: %v", err)
	}
}

// Create writes an empty setup the way `setup new` does, and it appears in the
// list.
func TestFyneSetupsCreateWritesAnEmptySetup(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	setups := fyneSetups{current: "default"}

	if err := setups.Create("streamer"); err != nil {
		t.Fatalf("Create: %v", err)
	}

	path := filepath.Join(home, "crower", "setups", "streamer", "crower.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}

	if !slices.Contains(setups.Names(), "streamer") {
		t.Errorf("Names() = %v, want the new setup", setups.Names())
	}

	if err := setups.Create("streamer"); err == nil {
		t.Error("Create accepted an existing setup")
	}

	if err := setups.Create("a/b"); err == nil {
		t.Error("Create accepted an invalid name")
	}
}

// Delete removes the folder, and refuses the default and the active setup.
func TestFyneSetupsDeleteRefusesTheDefaultAndTheActive(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	setups := fyneSetups{current: "work"}

	for _, name := range []string{"work", "old"} {
		if err := setups.Create(name); err != nil {
			t.Fatalf("Create %s: %v", name, err)
		}
	}

	if err := setups.Delete("default"); err == nil {
		t.Error("Delete removed the default setup")
	}

	if err := setups.Delete("work"); err == nil {
		t.Error("Delete removed the active setup")
	}

	if err := setups.Delete("missing"); err == nil {
		t.Error("Delete accepted a missing setup")
	}

	if err := setups.Delete("old"); err != nil {
		t.Fatalf("Delete old: %v", err)
	}

	if _, err := os.Stat(filepath.Join(home, "crower", "setups", "old")); !os.IsNotExist(err) {
		t.Errorf("the deleted setup's folder still exists: %v", err)
	}
}

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
