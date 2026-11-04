package catalog

import (
	"slices"
	"testing"
)

func TestNewKeepsTheGivenOrder(t *testing.T) {
	c := New(
		[]string{"zebra", "alpha"},
		[]Command{
			{Name: "alpha"},
			{Name: "middle"},
			{Name: "zebra"},
		},
	)

	want := []string{"zebra", "alpha", "middle"}
	if got := c.Names(); !slices.Equal(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

// An order that mentions a command the catalog does not hold, or mentions one
// twice, must not drop or duplicate anything.
func TestNewToleratesAStaleOrder(t *testing.T) {
	c := New(
		[]string{"gone", "b", "b", "a"},
		[]Command{{Name: "a"}, {Name: "b"}, {Name: "c"}},
	)

	want := []string{"b", "a", "c"}
	if got := c.Names(); !slices.Equal(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

func TestNewWithNoOrderIsSorted(t *testing.T) {
	c := New(nil, []Command{{Name: "c"}, {Name: "a"}, {Name: "b"}})

	if got, want := c.Names(), []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

// The alias owner must not depend on the display order: reordering the file
// must not change which command an alias executes.
func TestAliasOwnerIsIndependentOfDisplayOrder(t *testing.T) {
	commands := []Command{
		{Name: "b", Aliases: []string{"x"}},
		{Name: "a", Aliases: []string{"x"}},
	}

	for _, order := range [][]string{nil, {"b", "a"}, {"a", "b"}} {
		owner, ok := New(order, commands).AliasOwner("x")
		if !ok || owner != "a" {
			t.Errorf("order %v: alias x owned by %q (found %v), want a", order, owner, ok)
		}
	}
}

func TestAliasOwnerReportsUnclaimedAliases(t *testing.T) {
	c := New(nil, []Command{{Name: "a", Aliases: []string{"x"}}})

	if _, ok := c.AliasOwner("x"); !ok {
		t.Fatal("x is unclaimed, want a")
	}

	if _, ok := c.AliasOwner("nope"); ok {
		t.Fatal("nope is claimed, want unclaimed")
	}

	if _, ok := c.AliasOwner(""); ok {
		t.Fatal("the empty alias is claimed, want unclaimed")
	}
}

// A command whose alias equals its own name still claims it: that is what the
// loader has always done, and the warnings below decide what happens when
// another command wants the same key.
func TestSelfAliasClaimsTheKey(t *testing.T) {
	c := New(nil, []Command{{Name: "a", Aliases: []string{"a"}}})

	if owner, ok := c.AliasOwner("a"); !ok || owner != "a" {
		t.Errorf("alias a owned by %q (found %v), want a", owner, ok)
	}

	if warnings := c.Warnings(); len(warnings) != 0 {
		t.Errorf("Warnings() = %v, want none for a self alias", warnings)
	}
}

// The expected strings are the ones the loader printed before this package
// existed; users and their scripts read them.
func TestWarnings(t *testing.T) {
	tests := []struct {
		name     string
		commands []Command
		order    []string
		want     []string
	}{
		{
			name:     "no collision",
			commands: []Command{{Name: "a", Aliases: []string{"x"}}, {Name: "b", Aliases: []string{"y"}}},
		},
		{
			name:     "two commands claim the same alias",
			commands: []Command{{Name: "a", Aliases: []string{"x"}}, {Name: "b", Aliases: []string{"x"}}},
			want:     []string{`alias "x" is claimed by "a" and "b"; using "a"`},
		},
		{
			name:     "an alias is another command's name",
			commands: []Command{{Name: "a"}, {Name: "b", Aliases: []string{"a"}}},
			want:     []string{`alias "a" of "b" is also a command name; the name wins`},
		},
		{
			name:     "a contested alias that is also a name",
			commands: []Command{{Name: "a", Aliases: []string{"x"}}, {Name: "x", Aliases: []string{"x"}}},
			want: []string{
				`alias "x" is claimed by "a" and "x"`,
				`alias "x" of "a" is also a command name; the name wins`,
			},
		},
		{
			name:     "a duplicate alias inside one command warns once",
			commands: []Command{{Name: "a", Aliases: []string{"x"}}, {Name: "b", Aliases: []string{"x", "x"}}},
			want:     []string{`alias "x" is claimed by "a" and "b"; using "a"`},
		},
		{
			name: "three claimants",
			commands: []Command{
				{Name: "a", Aliases: []string{"x"}},
				{Name: "b", Aliases: []string{"x"}},
				{Name: "c", Aliases: []string{"x"}},
			},
			want: []string{
				`alias "x" is claimed by "a" and "b"; using "a"`,
				`alias "x" is claimed by "a" and "c"; using "a"`,
			},
		},
		{
			name:     "the warning order follows sorted names, not the file",
			commands: []Command{{Name: "b", Aliases: []string{"x"}}, {Name: "a", Aliases: []string{"x"}}},
			order:    []string{"b", "a"},
			want:     []string{`alias "x" is claimed by "a" and "b"; using "a"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.order, tt.commands).Warnings(); !slices.Equal(got, tt.want) {
				t.Errorf("Warnings() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Get and Commands hand out copies: a caller that writes through them must not
// be able to change what the catalog holds.
func TestAccessorsReturnCopies(t *testing.T) {
	c := New(nil, []Command{{Name: "a", Aliases: []string{"x"}, Steps: []string{"one"}}})

	got, ok := c.Get("a")
	if !ok {
		t.Fatal("Get(a) reported missing")
	}

	got.Aliases[0] = "changed"
	got.Steps[0] = "changed"

	stored, _ := c.Get("a")
	if stored.Aliases[0] != "x" || stored.Steps[0] != "one" {
		t.Errorf("the catalog changed through a copy: %+v", stored)
	}

	commands := c.Commands()
	commands[0].Name = "changed"

	if names := c.Names(); names[0] != "a" {
		t.Errorf("the catalog changed through Commands(): %v", names)
	}
}

func TestGetMissingCommand(t *testing.T) {
	if _, ok := New(nil, nil).Get("nope"); ok {
		t.Fatal("Get on an empty catalog found a command")
	}
}

func TestLen(t *testing.T) {
	c := New([]string{"a", "gone"}, []Command{{Name: "a"}, {Name: "b"}})

	if got := c.Len(); got != 2 {
		t.Errorf("Len() = %d, want 2", got)
	}
}

// A name wins over an alias. A key that is both must resolve to the command
// whose name it is, or that command would be unreachable by its own name --
// and `crower update` would silently change the wrong command.
func TestResolvePrefersTheName(t *testing.T) {
	c := New(nil, []Command{
		{Name: "a", Aliases: []string{"x"}},
		{Name: "x", Aliases: []string{"y"}},
	})

	got, found := c.Resolve("x")
	if !found {
		t.Fatal("Resolve(x) found nothing")
	}

	if got.Name != "x" {
		t.Errorf("Resolve(x) = %q, want the command named x", got.Name)
	}
}

func TestResolveFallsBackToAnAlias(t *testing.T) {
	c := New(nil, []Command{{Name: "greet", Aliases: []string{"hi"}}})

	got, found := c.Resolve("hi")
	if !found || got.Name != "greet" {
		t.Errorf("Resolve(hi) = %q (found %v), want greet", got.Name, found)
	}
}

func TestResolveUnknownKey(t *testing.T) {
	if _, found := New(nil, []Command{{Name: "greet"}}).Resolve("nope"); found {
		t.Fatal("Resolve(nope) found a command")
	}
}

// Metadata travels with the command and is copied, like every other slice.
// Shell and Cwd describe how a command runs, and an update that changes only
// its name must not drop them: the update starts from a copy of the old value.
func TestUpdateKeepsShellAndCwd(t *testing.T) {
	original := New(nil, []Command{{
		Name:  "deploy",
		Exec:  "echo hi",
		Shell: "pwsh",
		Cwd:   "~/projects/api",
	}})

	result, err := original.Update("deploy", Changes{Name: "ship"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if result.New.Shell != "pwsh" || result.New.Cwd != "~/projects/api" {
		t.Errorf("New = %+v, want the shell and cwd kept", result.New)
	}
}

func TestCommandCloneCopiesMetadata(t *testing.T) {
	original := New(nil, []Command{{
		Name:        "greet",
		Exec:        "echo hi",
		Description: "Say hi",
		Params:      []string{"name"},
	}})

	got, found := original.Get("greet")
	if !found || got.Description != "Say hi" || len(got.Params) != 1 || got.Params[0] != "name" {
		t.Fatalf("Get = %+v, %v; want the metadata", got, found)
	}

	got.Params[0] = "changed"

	again, _ := original.Get("greet")
	if again.Params[0] != "name" {
		t.Errorf("params = %v, want the catalog untouched", again.Params)
	}
}
