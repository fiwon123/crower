package session

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Reload keeps the same Session value but points it at another configuration,
// so every adapter built over it follows the switch.
func TestReloadPointsTheSessionAtAnotherConfiguration(t *testing.T) {
	first := newTestConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))
	if err := first.Commit(
		catalog.New(nil, []catalog.Command{{Name: "one", Exec: "echo one"}}),
		operation.Create, "one", "note",
	); err != nil {
		t.Fatalf("Commit one: %v", err)
	}

	secondPath := filepath.Join(t.TempDir(), "crower.yaml")
	second := newTestConfigAt(t, secondPath)
	if err := second.Commit(
		catalog.New(nil, []catalog.Command{{Name: "two", Exec: "echo two"}}),
		operation.Create, "two", "note",
	); err != nil {
		t.Fatalf("Commit two: %v", err)
	}

	if err := first.Reload(secondPath); err != nil {
		t.Fatalf("Reload: %v", err)
	}

	if _, found := first.Catalog().Get("one"); found {
		t.Error("the previous configuration is still loaded")
	}

	if _, found := first.Catalog().Get("two"); !found {
		t.Error("the reloaded configuration is missing its command")
	}

	if first.ConfigFile() != secondPath {
		t.Errorf("ConfigFile() = %q, want %q", first.ConfigFile(), secondPath)
	}

	if last := first.History.GetLast(); last == nil || last.CommandName != "two" {
		t.Errorf("history = %+v, want the other configuration's log", first.History.AllData)
	}
}

// A reload of an empty path is refused rather than leaving the session pointing
// at nothing.
func TestReloadRefusesAnEmptyPath(t *testing.T) {
	cfg := newTestConfig(t)

	if err := cfg.Reload(""); err == nil {
		t.Fatal("Reload accepted an empty path")
	}
}
