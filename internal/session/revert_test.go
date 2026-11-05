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
