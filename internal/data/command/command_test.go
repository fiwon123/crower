package command

import (
	"slices"
	"testing"
)

// AllAlias has no `omitempty`, so it is written to the configuration file on
// every save. A nil slice marshals as `null` in JSON and is dropped from a TOML
// table, so New normalises a nil list to an empty one instead of storing it.
func TestNewNormalisesNilAliases(t *testing.T) {
	c := New("greet", nil, "echo hello")

	if c.AllAlias == nil {
		t.Fatal("AllAlias is nil; the persisted form would be null")
	}

	if len(c.AllAlias) != 0 {
		t.Fatalf("AllAlias = %v, want empty", c.AllAlias)
	}
}

func TestNewStoresItsArguments(t *testing.T) {
	c := New("greet", []string{"hi", "hello"}, "echo hello")

	if c.Name != "greet" {
		t.Errorf("Name = %q, want greet", c.Name)
	}

	if !slices.Equal(c.AllAlias, []string{"hi", "hello"}) {
		t.Errorf("AllAlias = %v, want [hi hello]", c.AllAlias)
	}

	if c.Exec != "echo hello" {
		t.Errorf("Exec = %q, want echo hello", c.Exec)
	}

	if len(c.Steps) != 0 {
		t.Errorf("Steps = %v, want none", c.Steps)
	}
}

func TestMapDataAddAndGet(t *testing.T) {
	m := NewMapData()

	if m.Exists("greet") {
		t.Fatal("Exists reported a key that was never added")
	}

	m.Add("greet", New("greet", []string{"hi"}, "echo hello"))

	if !m.Exists("greet") {
		t.Fatal("Exists did not report the added key")
	}

	got := m.Get("greet")
	if got == nil {
		t.Fatal("Get returned nil for an added key")
	}

	if got.Name != "greet" || got.Exec != "echo hello" {
		t.Errorf("Get = %+v, want the added command", got)
	}
}

// The map is a map of values, so Get hands back a pointer to a copy. Writing a
// field through it must not change what the map holds -- that is why callers
// have to use Add or Update to persist a change.
func TestGetReturnsACopy(t *testing.T) {
	m := NewMapData()
	m.Add("greet", New("greet", []string{"hi"}, "echo hello"))

	got := m.Get("greet")
	got.Name = "changed"
	got.Exec = "changed"

	stored := m.Get("greet")
	if stored.Name != "greet" || stored.Exec != "echo hello" {
		t.Errorf("the map changed through a copy: %+v", stored)
	}
}

// Characterisation of a footgun the refactor has to remove: the struct is
// copied, but AllAlias is a slice, so the copy still shares the backing array
// with what the map holds. Mutating an element through Get therefore does
// write through. This test exists so a later change to copy semantics is a
// deliberate one, not an accident.
func TestGetSharesTheAliasBackingArray(t *testing.T) {
	m := NewMapData()
	m.Add("greet", New("greet", []string{"hi"}, "echo hello"))

	got := m.Get("greet")
	got.AllAlias[0] = "changed"

	if stored := m.Get("greet"); stored.AllAlias[0] != "changed" {
		t.Fatalf("AllAlias[0] = %q, want the mutation to be visible through the map", stored.AllAlias[0])
	}
}

func TestGetMissingValues(t *testing.T) {
	m := NewMapData()
	m.Add("greet", New("greet", nil, "echo hello"))

	tests := []struct {
		name string
		key  string
	}{
		{"missing key", "nope"},
		{"empty key", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.Get(tt.key); got != nil {
				t.Errorf("Get(%q) = %+v, want nil", tt.key, got)
			}
		})
	}
}

// Add accepts an empty key, but Get refuses to look one up. The entry is
// therefore stored and unreachable. Callers validate the key first, so this
// documents the asymmetry rather than guarding it.
func TestEmptyKeyIsStoredButUnreadable(t *testing.T) {
	m := NewMapData()
	m.Add("", New("", nil, "echo hello"))

	if !m.Exists("") {
		t.Fatal("Exists(\"\") = false, want the entry to be stored")
	}

	if got := m.Get(""); got != nil {
		t.Fatalf("Get(\"\") = %+v, want nil", got)
	}
}

func TestRemove(t *testing.T) {
	m := NewMapData()
	m.Add("greet", New("greet", nil, "echo hello"))

	if ok := m.Remove("greet"); !ok {
		t.Fatal("Remove of an existing key returned false")
	}

	if m.Exists("greet") {
		t.Fatal("the key survived Remove")
	}

	if ok := m.Remove("greet"); ok {
		t.Fatal("Remove of a missing key returned true")
	}
}

// Update writes the new value wholesale, so Steps travel with it: a caller that
// builds the new value without them replaces the stored steps with none.
func TestUpdateReplacesTheStoredValue(t *testing.T) {
	m := NewMapData()
	m.Add("greet", &Data{Name: "greet", AllAlias: []string{"hi"}, Exec: "echo old", Steps: []string{"one"}})

	ok := m.Update("greet", &Data{Name: "helloTen", AllAlias: []string{"hello"}, Exec: "echo new", Steps: []string{"one", "two"}})
	if !ok {
		t.Fatal("Update of an existing key returned false")
	}

	got := m.Get("greet")
	if got == nil {
		t.Fatal("the stored value disappeared")
	}

	if got.Name != "helloTen" || got.Exec != "echo new" {
		t.Errorf("Update left %+v, want the new name and exec", got)
	}

	if !slices.Equal(got.AllAlias, []string{"hello"}) {
		t.Errorf("AllAlias = %v, want [hello]", got.AllAlias)
	}

	if !slices.Equal(got.Steps, []string{"one", "two"}) {
		t.Errorf("Steps = %v, want the new steps", got.Steps)
	}
}

func TestUpdateOfAMissingKey(t *testing.T) {
	m := NewMapData()

	if ok := m.Update("nope", New("nope", nil, "echo hello")); ok {
		t.Fatal("Update of a missing key returned true")
	}

	if m.Exists("nope") {
		t.Fatal("Update created the missing entry")
	}
}
