package notes_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/shell"
)

// Each listing records the flag that reproduces it, and the positional values
// it was given.
func TestListNotesRecordReplayableCommands(t *testing.T) {
	tests := []struct {
		name string
		note string
		want string
	}{
		{"commands", notes.ListCommandsNote(), "crower_exec=crower list"},
		{"history", notes.ListHistoryNote(), "crower_exec=crower list --history"},
		{"system", notes.ListSystemNote(), "crower_exec=crower list --system"},
		{"syspath", notes.ListPathNote(), "crower_exec=crower list --syspath"},
		{"process", notes.ListProcessNote([]string{"go"}), "crower_exec=crower list --process " + shell.Quote("go")},
		{"folder", notes.ListFolderNote([]string{"~/projects"}), "crower_exec=crower list --folder " + shell.Quote("~/projects")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.HasPrefix(tt.note, "main_operation=list") {
				t.Errorf("note = %q, want the list operation first", tt.note)
			}

			if !strings.Contains(tt.note, tt.want) {
				t.Errorf("note = %q, want it to contain %q", tt.note, tt.want)
			}
		})
	}
}

// The note also names what was listed, so `list --history` reads as an audit
// line and not only as a replay recipe.
func TestListNotesRecordTheSubOperation(t *testing.T) {
	tests := []struct {
		note string
		want string
	}{
		{notes.ListCommandsNote(), "sub_operation=command"},
		{notes.ListProcessNote(nil), "sub_operation=process"},
		{notes.ListHistoryNote(), "sub_operation=history"},
		{notes.ListFolderNote(nil), "sub_operation=folder"},
		{notes.ListSystemNote(), "sub_operation=system"},
		{notes.ListPathNote(), "sub_operation=syspath"},
	}

	for _, tt := range tests {
		if !strings.Contains(tt.note, tt.want) {
			t.Errorf("note = %q, want it to contain %q", tt.note, tt.want)
		}
	}
}
