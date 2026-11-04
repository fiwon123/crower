package notes_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/shell"
)

// A deleted command is replayed by re-creating it, so the note has to carry the
// whole command and not only the name it was deleted by.
func TestDeleteCommandNoteRecordsTheWholeCommand(t *testing.T) {
	command := catalog.Command{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"}

	got := notes.DeleteCommandNote(command, []string{"greet"})

	want := "main_operation=delete;sub_operation=command;command_name=greet" +
		";command_alias=[hi];command_exec=echo hi;crower_exec=crower delete " + shell.Quote("greet")

	if got != want {
		t.Errorf("note = %q, want %q", got, want)
	}
}

func TestDeleteCommandNoteQuotesAnArgumentWithASpace(t *testing.T) {
	command := catalog.Command{Name: "greet", Exec: "echo hi"}

	got := notes.DeleteCommandNote(command, []string{"my file"})

	if want := "crower_exec=crower delete " + shell.Quote("my file"); !strings.Contains(got, want) {
		t.Errorf("note = %q, want it to contain %q", got, want)
	}
}

func TestDeleteLastNoteRecordsTheFlag(t *testing.T) {
	command := catalog.Command{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"}

	tests := []struct {
		name string
		op   operation.Operation
		want string
	}{
		{"create", operation.Create, "crower_exec=crower delete --create"},
		{"update", operation.Update, "crower_exec=crower delete --update"},
		{"execute", operation.Execute, "crower_exec=crower delete --execute"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := notes.DeleteLastNote(tt.op, command)

			if !strings.HasPrefix(got, "main_operation=delete;sub_operation=command") {
				t.Errorf("note = %q, want the delete/command fields first", got)
			}

			if !strings.Contains(got, "command_name=greet") {
				t.Errorf("note = %q, want the command name", got)
			}

			if !strings.Contains(got, tt.want) {
				t.Errorf("note = %q, want it to contain %q", got, tt.want)
			}
		})
	}
}

// An operation without a --last form records the command and no replay line,
// which is what it did before this package existed.
func TestDeleteLastNoteWithoutAFlagHasNoReplayLine(t *testing.T) {
	got := notes.DeleteLastNote(operation.Delete, catalog.Command{Name: "greet"})

	if strings.Contains(got, "crower_exec=") {
		t.Errorf("note = %q, want no crower_exec", got)
	}
}

func TestDeletePathNotes(t *testing.T) {
	tests := []struct {
		name string
		note func([]string) string
		sub  string
		exec string
	}{
		{"file", notes.DeleteFileNote, "sub_operation=file", "crower_exec=crower delete"},
		{"folder", notes.DeleteFolderNote, "sub_operation=folder", "crower_exec=crower delete"},
		{"system", notes.DeleteSystemVariableNote, "sub_operation=system", "crower_exec=crower delete --system"},
		{"syspath", notes.DeleteSystemPathNote, "sub_operation=syspath", "crower_exec=crower delete --syspath"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.note([]string{"notes.txt"})

			if !strings.HasPrefix(got, "main_operation=delete;"+tt.sub) {
				t.Errorf("note = %q, want it to start with main_operation=delete;%s", got, tt.sub)
			}

			want := tt.exec + " " + shell.Quote("notes.txt")
			if !strings.Contains(got, want) {
				t.Errorf("note = %q, want it to contain %q", got, want)
			}
		})
	}
}
