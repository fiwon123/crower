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

// The Fyne reverter goes through Session.RevertTo, so restoring an entry puts
// the command set back, drops the newer snapshots and records the revert.
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

	// Two entries back is the state after "a" was created.
	commands, history, err := newFyneReverter(cfg).Revert(2)
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
