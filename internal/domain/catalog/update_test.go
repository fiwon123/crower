package catalog

import (
	"errors"
	"slices"
	"testing"
)

func TestUpdateKeepsFieldsItWasNotGiven(t *testing.T) {
	original := New(nil, []Command{{
		Name:    "greet",
		Aliases: []string{"hi"},
		Exec:    "echo hi",
		Steps:   []string{"echo one", "echo two"},
	}})

	result, err := original.Update("greet", Changes{Exec: "echo new"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if result.New.Name != "greet" {
		t.Errorf("Name = %q, want greet", result.New.Name)
	}

	if !slices.Equal(result.New.Aliases, []string{"hi"}) {
		t.Errorf("Aliases = %v, want [hi]", result.New.Aliases)
	}

	if result.New.Exec != "echo new" {
		t.Errorf("Exec = %q, want echo new", result.New.Exec)
	}

	// update has no field for steps, so an update must carry them (#217).
	if !slices.Equal(result.New.Steps, []string{"echo one", "echo two"}) {
		t.Errorf("Steps = %v, want them preserved", result.New.Steps)
	}
}

// Update returns a new catalog and leaves the receiver alone, so a caller that
// reports a failure has not already changed what it was reading.
func TestUpdateDoesNotMutateTheReceiver(t *testing.T) {
	original := New(nil, []Command{{Name: "greet", Exec: "echo hi"}})

	if _, err := original.Update("greet", Changes{Exec: "echo new"}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	stored, _ := original.Get("greet")
	if stored.Exec != "echo hi" {
		t.Errorf("the receiver changed: Exec = %q", stored.Exec)
	}
}

func TestUpdateRenames(t *testing.T) {
	original := New(nil, []Command{{Name: "greet", Exec: "echo hi"}})

	result, err := original.Update("greet", Changes{Name: "hello"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if _, found := result.Catalog.Get("greet"); found {
		t.Error("the old name is still reachable")
	}

	if _, found := result.Catalog.Get("hello"); !found {
		t.Error("the new name is not reachable")
	}

	if result.Old.Name != "greet" || result.New.Name != "hello" {
		t.Errorf("Old/New = %q/%q, want greet/hello", result.Old.Name, result.New.Name)
	}
}

func TestUpdateKeepsThePositionAfterARename(t *testing.T) {
	original := New([]string{"first", "middle", "last"}, []Command{
		{Name: "first"}, {Name: "middle"}, {Name: "last"},
	})

	result, err := original.Update("middle", Changes{Name: "renamed"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	want := []string{"first", "renamed", "last"}
	if got := result.Catalog.Names(); !slices.Equal(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

func TestUpdateResolvesByAlias(t *testing.T) {
	original := New(nil, []Command{{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"}})

	result, err := original.Update("hi", Changes{Exec: "echo new"})
	if err != nil {
		t.Fatalf("Update by alias: %v", err)
	}

	if result.New.Name != "greet" {
		t.Errorf("updated %q, want greet", result.New.Name)
	}
}

func TestUpdateRefusesNothingToChange(t *testing.T) {
	original := New(nil, []Command{{Name: "greet", Exec: "echo hi"}})

	if _, err := original.Update("greet", Changes{}); !errors.Is(err, ErrNoChange) {
		t.Fatalf("Update error = %v, want %v", err, ErrNoChange)
	}
}

func TestUpdateRefusesUnknownKey(t *testing.T) {
	original := New(nil, []Command{{Name: "greet", Exec: "echo hi"}})

	if _, err := original.Update("nope", Changes{Exec: "x"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Update error = %v, want %v", err, ErrNotFound)
	}
}

func TestUpdateRefusesATakenName(t *testing.T) {
	original := New(nil, []Command{{Name: "a"}, {Name: "b"}})

	if _, err := original.Update("a", Changes{Name: "b"}); err == nil {
		t.Fatal("renaming onto another command's name was accepted")
	}
}

func TestUpdateRefusesANameThatIsAnAlias(t *testing.T) {
	original := New(nil, []Command{{Name: "a", Aliases: []string{"x"}}, {Name: "b"}})

	if _, err := original.Update("b", Changes{Name: "x"}); err == nil {
		t.Fatal("a name that is another command's alias was accepted")
	}
}

func TestUpdateRefusesANewAliasInUse(t *testing.T) {
	original := New(nil, []Command{{Name: "a", Aliases: []string{"x"}}, {Name: "b"}})

	if _, err := original.Update("b", Changes{Aliases: []string{"x"}}); err == nil {
		t.Fatal("an alias owned by another command was accepted")
	}
}

func TestUpdateRefusesANewAliasThatIsAName(t *testing.T) {
	original := New(nil, []Command{{Name: "a"}, {Name: "b"}})

	if _, err := original.Update("b", Changes{Aliases: []string{"a"}}); err == nil {
		t.Fatal("an alias that is another command's name was accepted")
	}
}

// A supplied list replaces the old one, so a claim this command does not own is
// reported and dropped rather than silently kept.
func TestUpdateDropsAForeignClaimFromASuppliedList(t *testing.T) {
	original := New(nil, []Command{
		{Name: "a", Aliases: []string{"x"}},
		{Name: "b", Aliases: []string{"x"}},
	})

	result, err := original.Update("b", Changes{Exec: "echo new", Aliases: []string{"x"}})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if len(result.New.Aliases) != 0 {
		t.Errorf("Aliases = %v, want the foreign claim dropped", result.New.Aliases)
	}

	want := []string{`alias "x" belongs to "a"; leaving it there`}
	if !slices.Equal(result.Warnings, want) {
		t.Errorf("Warnings = %v, want %v", result.Warnings, want)
	}
}

// An inherited list keeps a foreign claim dormant: the user changed another
// field, and dropping the alias would delete data they did not touch.
func TestUpdateKeepsAForeignClaimFromAnInheritedList(t *testing.T) {
	original := New(nil, []Command{
		{Name: "a", Aliases: []string{"x"}},
		{Name: "b", Aliases: []string{"x"}},
	})

	result, err := original.Update("b", Changes{Exec: "echo new"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if !slices.Equal(result.New.Aliases, []string{"x"}) {
		t.Errorf("Aliases = %v, want [x] kept dormant", result.New.Aliases)
	}

	if len(result.Warnings) != 0 {
		t.Errorf("Warnings = %v, want none for an inherited list", result.Warnings)
	}

	// The owner is unchanged: the alias still resolves to a.
	if owner, ok := result.Catalog.AliasOwner("x"); !ok || owner != "a" {
		t.Errorf("alias x owned by %q (found %v), want a", owner, ok)
	}
}

// A config may already carry a collision; an update of an unrelated field must
// still work on it.
func TestUpdateWorksOnAPreExistingCollision(t *testing.T) {
	original := New(nil, []Command{
		{Name: "a", Aliases: []string{"x"}},
		{Name: "b", Aliases: []string{"x"}},
	})

	if _, err := original.Update("a", Changes{Exec: "echo new"}); err != nil {
		t.Fatalf("Update of a colliding config: %v", err)
	}
}
