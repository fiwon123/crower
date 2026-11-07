package session

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// RevertTo replaces the in-memory catalog as well as the file. The Fyne window
// keeps one session alive across a revert; a stale catalog showed the commands
// the revert removed and wrote them back on the next save.
func TestRevertToReplacesTheInMemoryCatalog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crower.yaml")
	cfg := newTestConfigAt(t, path)

	first := catalog.New(nil, []catalog.Command{{Name: "a", Exec: "echo a"}})
	if err := cfg.Commit(first, operation.Create, "a", "note"); err != nil {
		t.Fatalf("Commit a: %v", err)
	}

	both := catalog.New(nil, []catalog.Command{
		{Name: "a", Exec: "echo a"},
		{Name: "b", Exec: "echo b"},
	})
	if err := cfg.Commit(both, operation.Create, "b", "note"); err != nil {
		t.Fatalf("Commit b: %v", err)
	}

	// Two entries back is the snapshot written by the first commit.
	entry, err := cfg.History.GetBeforeLast(2)
	if err != nil {
		t.Fatalf("GetBeforeLast: %v", err)
	}

	if err := cfg.RevertTo(entry); err != nil {
		t.Fatalf("RevertTo: %v", err)
	}

	if _, found := cfg.Catalog().Get("b"); found {
		t.Error("the in-memory catalog still holds the reverted command")
	}

	names := cfg.Catalog().Names()
	if len(names) != 1 || names[0] != "a" {
		t.Errorf("in-memory names = %v, want [a]", names)
	}

	// A later save must describe the restored state, not the one it replaced.
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	reloaded := newTestConfigAt(t, path)
	if _, found := reloaded.Catalog().Get("b"); found {
		t.Error("Save wrote the reverted command back to the configuration file")
	}

	if _, found := reloaded.Catalog().Get("a"); !found {
		t.Error("the restored command is missing from the configuration file")
	}
}

// RevertBack undoes the last operation: the delete's own snapshot still holds
// the deletion, so the state before it is what comes back.
func TestRevertBackUndoesOneOperation(t *testing.T) {
	cfg := newTestConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	withCommand := catalog.New(nil, []catalog.Command{{Name: "a", Exec: "echo a"}})
	if err := cfg.Commit(withCommand, operation.Create, "a", "note"); err != nil {
		t.Fatalf("Commit create: %v", err)
	}

	if err := cfg.Commit(catalog.New(nil, nil), operation.Delete, "a", "note"); err != nil {
		t.Fatalf("Commit delete: %v", err)
	}

	if _, err := cfg.RevertBack(1); err != nil {
		t.Fatalf("RevertBack: %v", err)
	}

	if _, found := cfg.Catalog().Get("a"); !found {
		t.Error("a did not come back from a revert")
	}
}

// Undoing every operation restores an empty set and drops the whole log.
func TestRevertBackUndoingEverythingLeavesAnEmptySet(t *testing.T) {
	cfg := newTestConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	for _, name := range []string{"a", "b"} {
		commands := catalog.New(nil, []catalog.Command{{Name: name, Exec: "echo " + name}})
		if err := cfg.Commit(commands, operation.Create, name, "note"); err != nil {
			t.Fatalf("Commit %s: %v", name, err)
		}
	}

	version, err := cfg.RevertBack(2)
	if err != nil {
		t.Fatalf("RevertBack: %v", err)
	}

	if version != 0 {
		t.Errorf("version = %d, want 0 for the empty set", version)
	}

	if names := cfg.Catalog().Names(); len(names) != 0 {
		t.Errorf("catalog = %v, want empty", names)
	}

	if len(cfg.History.AllData) != 0 {
		t.Errorf("history = %d entries, want none", len(cfg.History.AllData))
	}
}
