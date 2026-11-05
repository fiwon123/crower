package cmd

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage"
)

// The window shows the version the binary was built from, not a guess from the
// environment.
func TestFyneUpdaterCurrentIsTheBuildVersion(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	if got := newFyneUpdater(cfg, "v1.2.3").Current(); got != "v1.2.3" {
		t.Errorf("Current() = %q, want v1.2.3", got)
	}
}

// Restore previous uses the recorded version list, so a window can offer the
// same rollback the CLI has.
func TestFyneUpdaterReportsTheRecordedPreviousVersion(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	versions := storage.NewVersions(cfg.ConfigFolder())
	if err := versions.RecordInstalled("v1.2.2"); err != nil {
		t.Fatalf("RecordInstalled: %v", err)
	}

	previous, found := newFyneUpdater(cfg, "v1.2.3").Previous()
	if !found || previous != "v1.2.2" {
		t.Fatalf("Previous() = %q, %v; want v1.2.2, true", previous, found)
	}
}

// A fresh installation has nothing to restore, and the window must say so
// rather than offer a button that cannot work.
func TestFyneUpdaterHasNoPreviousVersionOnAFreshInstallation(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	if previous, found := newFyneUpdater(cfg, "v1.2.3").Previous(); found {
		t.Errorf("Previous() = %q, true; want no version", previous)
	}
}

// Restart spawns the replacement exactly once and never runs it in a test.
func TestFyneUpdaterRestartSpawnsTheReplacement(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	spawned := 0
	updater := newFyneUpdater(cfg, "v1.2.3")
	updater.start = func() error {
		spawned++
		return nil
	}

	if err := updater.Restart(); err != nil {
		t.Fatalf("Restart: %v", err)
	}

	if spawned != 1 {
		t.Errorf("spawned %d times, want 1", spawned)
	}
}
