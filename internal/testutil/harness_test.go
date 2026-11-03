package testutil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/testutil"
)

func TestNewCreatesAnIsolatedInstallation(t *testing.T) {
	cfg, paths, err := testutil.New(t)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if cfg == nil {
		t.Fatal("New returned no configuration")
	}

	if len(paths) != 3 {
		t.Fatalf("New returned %d folders, want 3", len(paths))
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("fixture folder %s: %v", path, err)
			continue
		}

		if !info.IsDir() {
			t.Errorf("%s is not a directory", path)
		}
	}

	// The configuration lives inside the fixture, not next to the caller's home
	// directory: a test must never touch the machine it runs on.
	root := filepath.Dir(paths[0])

	if got := filepath.Dir(cfg.ConfigFile()); got != root {
		t.Errorf("config folder = %q, want %q", got, root)
	}
}

// Every call gets its own directory, so two tests -- or two calls in one test --
// cannot see each other's commands.
func TestNewGivesEachCallItsOwnDirectory(t *testing.T) {
	first, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("first New: %v", err)
	}

	second, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("second New: %v", err)
	}

	if first.ConfigFolder() == second.ConfigFolder() {
		t.Fatalf("both fixtures share %q", first.ConfigFolder())
	}
}
