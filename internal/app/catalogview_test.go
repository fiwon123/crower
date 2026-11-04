package app_test

import (
	"testing"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/testutil"
)

// Two commands can declare the same alias -- a hand-edited or legacy config
// does. The catalog gives it to the first declarer in sorted order, and
// ReplaceCatalog is what the delete use case commits through, so this pins the
// claim that the alias moves to the other declarer at once rather than going
// dormant until a reload.
func TestReplaceCatalogHandsAContestedAliasToTheOtherDeclarer(t *testing.T) {
	cfg, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("testutil.New: %v", err)
	}

	cfg.AllCommandsByName.Add("a", command.New("a", []string{"x"}, "echo a"))
	cfg.AllCommandsByName.Add("b", command.New("b", []string{"x"}, "echo b"))

	// Build the maps the loader would have built from that config.
	cfg.ReplaceCatalog(cfg.Catalog())

	if owner := cfg.AllCommandsByAlias.Get("x"); owner == nil || owner.Name != "a" {
		t.Fatalf("alias x = %+v, want a before the delete", owner)
	}

	remaining, _, ok := cfg.Catalog().Remove("a")
	if !ok {
		t.Fatal("Remove(a) = false")
	}

	cfg.ReplaceCatalog(remaining)

	if owner := cfg.AllCommandsByAlias.Get("x"); owner == nil || owner.Name != "b" {
		t.Errorf("alias x = %+v, want b after deleting its owner", owner)
	}
}

// Deleting the declarer that did not own the alias must leave it where it is.
func TestReplaceCatalogKeepsAnAliasWhenItsNonOwnerIsDeleted(t *testing.T) {
	cfg, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("testutil.New: %v", err)
	}

	cfg.AllCommandsByName.Add("a", command.New("a", []string{"x"}, "echo a"))
	cfg.AllCommandsByName.Add("b", command.New("b", []string{"x"}, "echo b"))

	cfg.ReplaceCatalog(cfg.Catalog())

	remaining, _, ok := cfg.Catalog().Remove("b")
	if !ok {
		t.Fatal("Remove(b) = false")
	}

	cfg.ReplaceCatalog(remaining)

	if owner := cfg.AllCommandsByAlias.Get("x"); owner == nil || owner.Name != "a" {
		t.Errorf("alias x = %+v, want a after deleting the other declarer", owner)
	}
}
