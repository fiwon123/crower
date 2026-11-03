package updatedom

import (
	"strings"
	"testing"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/shell"
)

func TestUpdateNoteRecordsTheFlags(t *testing.T) {
	quote := shell.Quote

	oldCommand := command.New("old", []string{"o"}, "echo old")
	newCommand := command.New("new", []string{"n"}, "echo new")

	built := newUpdateCommmandNote(oldCommand, newCommand)

	want := "crower_exec=crower update --name " + quote("new") +
		" --alias " + quote("n") +
		" --exec " + quote("echo new") +
		" " + quote("new")

	if !strings.Contains(built, want) {
		t.Errorf("note = %q, want it to contain %q", built, want)
	}
}

// Every field the prompt would ask for is recorded, not only the one that
// changed, so replay does not stop to ask for the others.
func TestUpdateNoteRecordsEveryField(t *testing.T) {
	quote := shell.Quote

	oldCommand := command.New("greet", []string{"g"}, "echo old")
	newCommand := command.New("greet", []string{"g"}, "echo new")

	built := newUpdateCommmandNote(oldCommand, newCommand)

	want := "crower_exec=crower update --name " + quote("greet") +
		" --alias " + quote("g") +
		" --exec " + quote("echo new") +
		" " + quote("greet")
	if !strings.Contains(built, want) {
		t.Errorf("note = %q, want it to contain %q", built, want)
	}
}

func TestUpdateNoteEncodesAliasesAsCSV(t *testing.T) {
	oldCommand := command.New("greet", nil, "echo hi")
	newCommand := command.New("greet", []string{"a,b", "g2"}, "echo hi")

	built := newUpdateCommmandNote(oldCommand, newCommand)

	want := "crower_exec=crower update --name " + shell.Quote("greet") +
		" --alias " + shell.Quote(`"a,b",g2`) +
		" --exec " + shell.Quote("echo hi") +
		" " + shell.Quote("greet")
	if !strings.Contains(built, want) {
		t.Errorf("note = %q, want it to contain %q", built, want)
	}
}

func TestUpdateLastNoteRecordsTheFlags(t *testing.T) {
	oldCommand := command.New("greet", nil, "echo old")
	newCommand := command.New("greet", nil, "echo new")

	built := newUpdateLastNote(operation.Create, oldCommand, newCommand)

	want := "crower_exec=crower update --create --name " + shell.Quote("greet") +
		" --exec " + shell.Quote("echo new")
	if !strings.Contains(built, want) {
		t.Errorf("note = %q, want it to contain %q", built, want)
	}
}

// Comparing aliases must not sort the slices stored on the command map.
func TestHasDiffAliasesDoesNotMutate(t *testing.T) {
	oldAliases := []string{"b", "a"}
	newAliases := []string{"a", "b"}

	if hasDiffAlieses(oldAliases, newAliases) {
		t.Error("the same aliases in another order were reported as different")
	}

	if strings.Join(oldAliases, ",") != "b,a" || strings.Join(newAliases, ",") != "a,b" {
		t.Errorf("inputs were mutated: old=%v new=%v", oldAliases, newAliases)
	}
}
