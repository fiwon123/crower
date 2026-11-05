package cmd

import (
	"io"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/session"
)

func newFyneEditorConfig(t *testing.T) *session.Session {
	t.Helper()

	cfg, err := session.Load(
		filepath.Join(t.TempDir(), "crower.yaml"),
		session.Options{Streams: session.Streams{
			In:  strings.NewReader(""),
			Out: io.Discard,
			Err: io.Discard,
		}},
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	return cfg
}

// The Fyne editor goes through catalog.Update and Session.Commit, so a rename,
// a changed exec and an alias change are persisted and recorded like the CLI.
func TestFyneEditorUpdateCommitsTheChange(t *testing.T) {
	cfg, err := session.Load(
		filepath.Join(t.TempDir(), "crower.yaml"),
		session.Options{Streams: session.Streams{
			In:  strings.NewReader(""),
			Out: io.Discard,
			Err: io.Discard,
		}},
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	cfg.ReplaceCatalog(catalog.New(nil, []catalog.Command{{
		Name:    "greet",
		Aliases: []string{"hi"},
		Exec:    "echo hi",
	}}))

	commands, history, err := newFyneEditor(cfg).Update("greet", catalog.Changes{
		Name:         "ship",
		Aliases:      []string{"s"},
		Exec:         "echo ship",
		ClearAliases: false,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, found := commands.Get("ship")
	if !found || got.Exec != "echo ship" {
		t.Fatalf("commands = %+v, %v; want ship with the new exec", got, found)
	}

	if len(history.AllData) != 1 || history.AllData[0].CommandName != "ship" {
		t.Errorf("history = %+v, want one update entry for ship", history.AllData)
	}

	if stored, found := cfg.Catalog().Get("ship"); !found || stored.Exec != "echo ship" {
		t.Errorf("session catalog = %+v, %v; want the change persisted", stored, found)
	}
}

// Create goes through catalog.Add and records a create entry.
func TestFyneEditorCreateCommitsTheCommand(t *testing.T) {
	cfg := newFyneEditorConfig(t)

	commands, history, err := newFyneEditor(cfg).Create(catalog.Command{Name: "ship", Aliases: []string{"s"}, Exec: "echo ship"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, found := commands.Get("ship")
	if !found || got.Exec != "echo ship" || !slices.Equal(got.Aliases, []string{"s"}) {
		t.Fatalf("commands = %+v, %v; want ship with its alias and exec", got, found)
	}

	if len(history.AllData) != 1 || history.AllData[0].CommandName != "ship" {
		t.Errorf("history = %+v, want one create entry for ship", history.AllData)
	}
}

// Delete goes through catalog.Remove and records the command it removed.
func TestFyneEditorDeleteCommitsTheChange(t *testing.T) {
	cfg := newFyneEditorConfig(t)
	if _, _, err := newFyneEditor(cfg).Create(catalog.Command{Name: "ship", Exec: "echo ship"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	commands, history, err := newFyneEditor(cfg).Delete("ship")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if _, found := commands.Get("ship"); found {
		t.Error("Delete left the command in the set")
	}

	if len(history.AllData) != 2 || history.AllData[1].Operation != operation.Delete {
		t.Errorf("history = %+v, want a delete entry", history.AllData)
	}
}

// A create with a name that already exists is refused and writes nothing.
func TestFyneEditorCreateRefusesACollision(t *testing.T) {
	cfg := newFyneEditorConfig(t)
	if _, _, err := newFyneEditor(cfg).Create(catalog.Command{Name: "ship", Exec: "echo ship"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	before := len(cfg.History.AllData)

	if _, _, err := newFyneEditor(cfg).Create(catalog.Command{Name: "ship", Exec: "echo again"}); err == nil {
		t.Fatal("Create accepted a duplicate name")
	}

	if got := len(cfg.History.AllData); got != before {
		t.Errorf("history grew from %d to %d for a refused create", before, got)
	}
}

// A collision is the catalog's error, and nothing is written.
func TestFyneEditorUpdateRefusesACollision(t *testing.T) {
	cfg, err := session.Load(
		filepath.Join(t.TempDir(), "crower.yaml"),
		session.Options{Streams: session.Streams{
			In:  strings.NewReader(""),
			Out: io.Discard,
			Err: io.Discard,
		}},
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	cfg.ReplaceCatalog(catalog.New(nil, []catalog.Command{
		{Name: "greet", Exec: "echo hi"},
		{Name: "build", Exec: "make"},
	}))

	_, _, err = newFyneEditor(cfg).Update("greet", catalog.Changes{Name: "build"})
	if err == nil {
		t.Fatal("Update accepted a name that already exists")
	}

	if _, found := cfg.Catalog().Get("build"); !found {
		t.Fatal("the refused update removed the existing command")
	}

	if len(cfg.History.AllData) != 0 {
		t.Errorf("history = %+v, want no entry for a refused update", cfg.History.AllData)
	}
}
