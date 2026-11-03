package notes_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/shell"
)

// The positional form is a replay line, so it records every field the prompt
// would ask for -- not only the one that changed -- or replay stops to ask for
// the rest. The key is the command's final name.
func TestUpdateNoteRecordsTheFlags(t *testing.T) {
	before := catalog.Command{Name: "old", Aliases: []string{"o"}, Exec: "echo old"}
	after := catalog.Command{Name: "new", Aliases: []string{"n"}, Exec: "echo new"}

	got := notes.UpdateNote(before, after, "")

	want := "crower_exec=crower update --name " + shell.Quote("new") +
		" --alias " + shell.Quote("n") +
		" --exec " + shell.Quote("echo new") +
		" " + shell.Quote("new")

	if !strings.Contains(got, want) {
		t.Errorf("note = %q, want it to contain %q", got, want)
	}
}

func TestUpdateNoteRecordsEveryFieldEvenWhenOnlyOneChanged(t *testing.T) {
	before := catalog.Command{Name: "greet", Aliases: []string{"g"}, Exec: "echo old"}
	after := catalog.Command{Name: "greet", Aliases: []string{"g"}, Exec: "echo new"}

	got := notes.UpdateNote(before, after, "")

	want := "crower_exec=crower update --name " + shell.Quote("greet") +
		" --alias " + shell.Quote("g") +
		" --exec " + shell.Quote("echo new") +
		" " + shell.Quote("greet")

	if !strings.Contains(got, want) {
		t.Errorf("note = %q, want it to contain %q", got, want)
	}
}

// Aliases travel as CSV, so one containing a comma survives replay.
func TestUpdateNoteEncodesAliasesAsCSV(t *testing.T) {
	before := catalog.Command{Name: "greet", Exec: "echo hi"}
	after := catalog.Command{Name: "greet", Aliases: []string{"a,b", "g2"}, Exec: "echo hi"}

	got := notes.UpdateNote(before, after, "")

	want := "crower_exec=crower update --name " + shell.Quote("greet") +
		" --alias " + shell.Quote(`"a,b",g2`) +
		" --exec " + shell.Quote("echo hi") +
		" " + shell.Quote("greet")

	if !strings.Contains(got, want) {
		t.Errorf("note = %q, want it to contain %q", got, want)
	}
}

// A --last form carries no positional key: the flags apply to the command the
// last operation touched.
func TestUpdateNoteLastForms(t *testing.T) {
	before := catalog.Command{Name: "greet", Exec: "echo old"}
	after := catalog.Command{Name: "greet", Exec: "echo new"}

	tests := []struct {
		last operation.Operation
		want string
	}{
		{operation.Update, "--last"},
		{operation.Create, "--create"},
		{operation.Execute, "--execute"},
	}

	for _, tt := range tests {
		t.Run(string(tt.last), func(t *testing.T) {
			got := notes.UpdateNote(before, after, tt.last)

			want := "crower_exec=crower update " + tt.want + " --name " + shell.Quote("greet") +
				" --exec " + shell.Quote("echo new")

			if !strings.Contains(got, want) {
				t.Errorf("note = %q, want it to contain %q", got, want)
			}

			if strings.Contains(got, shell.Quote("greet")+"\n") {
				t.Errorf("a --last note carries a positional key: %q", got)
			}
		})
	}
}

// The diff records the pair of values that differ, and nothing for a field that
// did not change.
func TestUpdateNoteRecordsTheDiff(t *testing.T) {
	before := catalog.Command{Name: "old", Exec: "echo old"}
	after := catalog.Command{Name: "new", Exec: "echo old"}

	got := notes.UpdateNote(before, after, "")

	for _, want := range []string{"old_command_name=old", "command_name=new"} {
		if !strings.Contains(got, want) {
			t.Errorf("note = %q, want it to contain %q", got, want)
		}
	}

	if strings.Contains(got, "command_exec=") {
		t.Errorf("note records an unchanged exec: %q", got)
	}
}

// Comparing aliases must not reorder the slices the command holds.
func TestUpdateNoteDoesNotMutateTheAliasSlices(t *testing.T) {
	beforeAliases := []string{"b", "a"}
	afterAliases := []string{"a", "b"}

	notes.UpdateNote(
		catalog.Command{Name: "greet", Aliases: beforeAliases},
		catalog.Command{Name: "greet", Aliases: afterAliases},
		"",
	)

	if strings.Join(beforeAliases, ",") != "b,a" || strings.Join(afterAliases, ",") != "a,b" {
		t.Errorf("inputs were mutated: before=%v after=%v", beforeAliases, afterAliases)
	}
}
