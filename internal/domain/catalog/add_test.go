package catalog

import (
	"slices"
	"testing"
)

func TestAddAppendsTheCommand(t *testing.T) {
	original := New(nil, []Command{{Name: "alpha", Exec: "echo a"}})

	result, err := original.Add(Command{Name: "beta", Aliases: []string{"b"}, Exec: "echo b"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if got, want := result.Names(), []string{"alpha", "beta"}; !slices.Equal(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}

	if owner, found := result.AliasOwner("b"); !found || owner != "beta" {
		t.Errorf("alias b = %q, %v; want beta, true", owner, found)
	}

	// The receiver is unchanged: Add returns a new catalog.
	if original.Len() != 1 {
		t.Errorf("the original catalog now holds %d commands, want 1", original.Len())
	}
}

func TestAddKeepsTheDisplayOrder(t *testing.T) {
	original := New([]string{"zebra", "alpha"}, []Command{
		{Name: "alpha", Exec: "echo a"},
		{Name: "zebra", Exec: "echo z"},
	})

	result, err := original.Add(Command{Name: "middle", Exec: "echo m"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if got, want := result.Names(), []string{"zebra", "alpha", "middle"}; !slices.Equal(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

func TestAddRefusesWhatCannotBeAdded(t *testing.T) {
	base := New(nil, []Command{
		{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"},
	})

	tests := []struct {
		name    string
		command Command
		want    string
	}{
		{"empty name", Command{Exec: "echo x"}, "empty name"},
		{"empty exec", Command{Name: "x"}, "empty exec"},
		{"name already used", Command{Name: "greet", Exec: "echo x"}, "found name, command already added"},
		{"name is another command's alias", Command{Name: "hi", Exec: "echo x"}, `"hi" is already an alias of "greet"`},
		{"alias already used", Command{Name: "x", Aliases: []string{"hi"}, Exec: "echo x"}, "found alias, command already added"},
		{"alias is another command's name", Command{Name: "x", Aliases: []string{"greet"}, Exec: "echo x"}, "found alias, command already added"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := base.Add(tt.command)
			if err == nil {
				t.Fatalf("Add(%+v) returned no error", tt.command)
			}

			if err.Error() != tt.want {
				t.Errorf("error = %q, want %q", err, tt.want)
			}

			if result.Len() != base.Len() {
				t.Errorf("a refused Add changed the catalog to %d commands", result.Len())
			}
		})
	}
}

// A command whose own name appears in its alias list is accepted, as before:
// the name is not yet in the catalog when the aliases are checked, and the
// name index answers first anyway.
func TestAddAcceptsACommandThatListsItsOwnNameAsAnAlias(t *testing.T) {
	result, err := New(nil, nil).Add(Command{Name: "greet", Aliases: []string{"greet"}, Exec: "echo hi"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if command, found := result.Resolve("greet"); !found || command.Name != "greet" {
		t.Errorf("Resolve(greet) = %+v, %v; want the command", command, found)
	}
}

// An empty alias is ignored rather than refused: it names nothing.
func TestAddIgnoresAnEmptyAlias(t *testing.T) {
	result, err := New(nil, nil).Add(Command{Name: "greet", Aliases: []string{"", "hi"}, Exec: "echo hi"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if owner, found := result.AliasOwner("hi"); !found || owner != "greet" {
		t.Errorf("alias hi = %q, %v; want greet, true", owner, found)
	}
}
