package catalog

import (
	"slices"
	"testing"
)

func removeFixture() Catalog {
	return New(
		[]string{"beta", "alpha", "gamma"},
		[]Command{
			{Name: "alpha", Aliases: []string{"a"}, Exec: "echo a"},
			{Name: "beta", Aliases: []string{"b"}, Exec: "echo b"},
			{Name: "gamma", Exec: "echo g"},
		},
	)
}

func TestRemoveByCommandName(t *testing.T) {
	result, removed, ok := removeFixture().Remove("beta")
	if !ok {
		t.Fatal("Remove(beta) = false, want true")
	}

	if removed.Name != "beta" || removed.Exec != "echo b" {
		t.Errorf("removed = %+v, want the beta command", removed)
	}

	if got, want := result.Names(), []string{"alpha", "gamma"}; !slices.Equal(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}

	if _, found := result.Get("beta"); found {
		t.Error("beta is still in the catalog")
	}

	if _, found := result.AliasOwner("b"); found {
		t.Error("alias b still resolves")
	}
}

// An alias resolves to the command that owns it, and Remove reports that
// command -- not the key the user typed.
func TestRemoveByAlias(t *testing.T) {
	result, removed, ok := removeFixture().Remove("a")
	if !ok {
		t.Fatal("Remove(a) = false, want true")
	}

	if removed.Name != "alpha" {
		t.Errorf("removed.Name = %q, want alpha", removed.Name)
	}

	if _, found := result.Get("alpha"); found {
		t.Error("alpha is still in the catalog")
	}
}

func TestRemoveUnknownKeyChangesNothing(t *testing.T) {
	original := removeFixture()

	result, removed, ok := original.Remove("nothing")
	if ok {
		t.Fatal("Remove(nothing) = true, want false")
	}

	if removed.Name != "" || removed.Exec != "" || len(removed.Aliases) != 0 {
		t.Errorf("removed = %+v, want the zero command", removed)
	}

	if !slices.Equal(result.Names(), original.Names()) {
		t.Errorf("Names() = %v, want the catalog unchanged", result.Names())
	}

	if owner, found := result.AliasOwner("a"); !found || owner != "alpha" {
		t.Errorf("alias a = %q, %v; want alpha, true", owner, found)
	}
}

// A legacy config can hold two commands that declare the same alias. The first
// claim in sorted order wins, and deleting the winner hands the alias to the
// other declarer -- which is what reloading the written file does. The legacy
// handler instead dropped the index entry and left the alias dormant until a
// restart, so the running process and the next one disagreed.
func TestRemoveHandsAContestedAliasToItsOtherDeclarer(t *testing.T) {
	c := New(nil, []Command{
		{Name: "a", Aliases: []string{"x"}, Exec: "echo a"},
		{Name: "b", Aliases: []string{"x"}, Exec: "echo b"},
	})

	if owner, _ := c.AliasOwner("x"); owner != "a" {
		t.Fatalf("alias x resolves to %q, want a", owner)
	}

	result, _, ok := c.Remove("a")
	if !ok {
		t.Fatal("Remove(a) = false, want true")
	}

	owner, found := result.AliasOwner("x")
	if !found || owner != "b" {
		t.Errorf("alias x = %q, %v; want b, true", owner, found)
	}

	if got, found := result.Resolve("x"); !found || got.Name != "b" {
		t.Errorf("Resolve(x) = %+v, %v; want b, true", got, found)
	}
}

// The removed command's slices are returned as a copy, so a caller cannot
// change the catalog it came from.
func TestRemoveReturnsACopy(t *testing.T) {
	original := removeFixture()

	_, removed, ok := original.Remove("alpha")
	if !ok {
		t.Fatal("Remove(alpha) = false, want true")
	}

	removed.Aliases[0] = "changed"

	if got, _ := original.Get("alpha"); got.Aliases[0] != "a" {
		t.Errorf("catalog alias = %q, want the original a", got.Aliases[0])
	}
}
