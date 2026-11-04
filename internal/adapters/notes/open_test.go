package notes_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/shell"
)

func TestOpenNotesRecordReplayableCommands(t *testing.T) {
	tests := []struct {
		name string
		note string
		want string
		sub  string
	}{
		{"paths", notes.OpenNote([]string{"a.txt", "b.txt"}), "crower_exec=crower open " + shell.Quote("a.txt") + " " + shell.Quote("b.txt"), ""},
		{"file", notes.OpenFileNote([]string{"notes.txt"}), "crower_exec=crower open --file " + shell.Quote("notes.txt"), "sub_operation=file"},
		{"folder", notes.OpenFolderNote([]string{"~/p"}), "crower_exec=crower open --folder " + shell.Quote("~/p"), "sub_operation=folder"},
		{"system", notes.OpenSystemNote(), "crower_exec=crower open --system", "sub_operation=system"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.note, tt.want) {
				t.Errorf("note = %q, want it to contain %q", tt.note, tt.want)
			}

			if tt.sub != "" && !strings.Contains(tt.note, tt.sub) {
				t.Errorf("note = %q, want it to contain %q", tt.note, tt.sub)
			}
		})
	}
}

// The positional form has no sub-operation; it is what it always was.
func TestOpenNoteHasNoSubOperation(t *testing.T) {
	if got := notes.OpenNote(nil); strings.Contains(got, "sub_operation") {
		t.Errorf("note = %q, want no sub-operation", got)
	}
}
