package session

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage/command"
)

// newCatalogSession builds an empty run whose command map a test can seed.
func newCatalogSession(t *testing.T) *Session {
	t.Helper()

	cfg, err := New(
		filepath.Join(t.TempDir(), "crower.yaml"),
		nil,
		command.NewMapData(),
		command.NewMapData(),
		Options{Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard}},
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	return cfg
}

// Two commands can declare the same alias -- a hand-edited or legacy config
// does. The catalog gives it to the first declarer in sorted order, and
// ReplaceCatalog is what the delete use case commits through, so this pins the
// claim that the alias moves to the other declarer at once rather than going
// dormant until a reload.
func TestReplaceCatalogHandsAContestedAliasToTheOtherDeclarer(t *testing.T) {
	cfg := newCatalogSession(t)

	cfg.allCommandsByName.Add("a", command.New("a", []string{"x"}, "echo a"))
	cfg.allCommandsByName.Add("b", command.New("b", []string{"x"}, "echo b"))

	// Build the maps the loader would have built from that config.
	cfg.ReplaceCatalog(cfg.Catalog())

	if owner := cfg.allCommandsByAlias.Get("x"); owner == nil || owner.Name != "a" {
		t.Fatalf("alias x = %+v, want a before the delete", owner)
	}

	remaining, _, ok := cfg.Catalog().Remove("a")
	if !ok {
		t.Fatal("Remove(a) = false")
	}

	cfg.ReplaceCatalog(remaining)

	if owner := cfg.allCommandsByAlias.Get("x"); owner == nil || owner.Name != "b" {
		t.Errorf("alias x = %+v, want b after deleting its owner", owner)
	}
}

// Deleting the declarer that did not own the alias must leave it where it is.
func TestReplaceCatalogKeepsAnAliasWhenItsNonOwnerIsDeleted(t *testing.T) {
	cfg := newCatalogSession(t)

	cfg.allCommandsByName.Add("a", command.New("a", []string{"x"}, "echo a"))
	cfg.allCommandsByName.Add("b", command.New("b", []string{"x"}, "echo b"))

	cfg.ReplaceCatalog(cfg.Catalog())

	remaining, _, ok := cfg.Catalog().Remove("b")
	if !ok {
		t.Fatal("Remove(b) = false")
	}

	cfg.ReplaceCatalog(remaining)

	if owner := cfg.allCommandsByAlias.Get("x"); owner == nil || owner.Name != "a" {
		t.Errorf("alias x = %+v, want a after deleting the other declarer", owner)
	}
}
