package cmd

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/session"
)

func loadFyneConfigAt(t *testing.T, path string) *session.Session {
	t.Helper()

	cfg, err := session.Load(path, session.Options{Streams: session.Streams{
		In:  strings.NewReader(""),
		Out: io.Discard,
		Err: io.Discard,
	}})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	return cfg
}

// One step undoes one operation: the snapshot of the entry before it is the
// state the user had (#483).
func TestFyneReverterRestoresAnEarlierState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crower.yaml")
	cfg := loadFyneConfigAt(t, path)

	editor := newFyneEditor(cfg)
	if _, _, err := editor.Create(catalog.Command{Name: "a", Exec: "echo a"}); err != nil {
		t.Fatalf("Create a: %v", err)
	}

	if _, _, err := editor.Create(catalog.Command{Name: "b", Exec: "echo b"}); err != nil {
		t.Fatalf("Create b: %v", err)
	}

	// Undoing the last operation lands on the state after "a" was created.
	commands, history, err := newFyneReverter(cfg).Revert(1)
	if err != nil {
		t.Fatalf("Revert: %v", err)
	}

	if _, found := commands.Get("b"); found {
		t.Error("b survived a revert that should have removed it")
	}

	if _, found := commands.Get("a"); !found {
		t.Error("a is missing after the revert")
	}

	last := history.GetLast()
	if last == nil || last.Operation != operation.Revert {
		t.Errorf("last history entry = %+v, want a revert entry", last)
	}

	reloaded := loadFyneConfigAt(t, path)
	if _, found := reloaded.Catalog().Get("b"); found {
		t.Error("the reverted set was not written to the configuration file")
	}

	if _, found := reloaded.Catalog().Get("a"); !found {
		t.Error("the restored command is missing from the configuration file")
	}
}

// Restoring after a delete brings the command back: that is the report that
// exposed the off-by-one — the delete entry's own snapshot already lacks it.
func TestFyneReverterBringsADeletedCommandBack(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))
	editor := newFyneEditor(cfg)

	if _, _, err := editor.Create(catalog.Command{Name: "greet", Exec: "echo hi"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if _, _, err := editor.Delete("greet"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	commands, _, err := newFyneReverter(cfg).Revert(1)
	if err != nil {
		t.Fatalf("Revert: %v", err)
	}

	if _, found := commands.Get("greet"); !found {
		t.Error("restoring after a delete did not bring the command back")
	}
}

// A zero step count would index past the end of the log, so the port refuses it
// rather than panicking.
func TestFyneReverterRefusesZeroSteps(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	if _, _, err := newFyneReverter(cfg).Revert(0); err == nil {
		t.Fatal("Revert(0) was accepted")
	}
}

// A step count larger than the log is reported, not applied.
func TestFyneReverterRefusesMoreStepsThanHistoryHolds(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	if _, _, err := newFyneReverter(cfg).Revert(1); err == nil {
		t.Fatal("Revert(1) was accepted over an empty history")
	}
}
