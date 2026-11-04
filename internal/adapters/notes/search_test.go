package notes_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/shell"
)

func TestSearchNotesRecordReplayableCommands(t *testing.T) {
	tests := []struct {
		name string
		note string
		want string
	}{
		{"browser", notes.SearchBrowserNote([]string{"golang context"}), "crower_exec=crower search --browser " + shell.Quote("golang context")},
		{"file", notes.SearchFileNote([]string{"a.txt", "."}), "crower_exec=crower search --file " + shell.Quote("a.txt") + " " + shell.Quote(".")},
		{"folder", notes.SearchFolderNote([]string{"logs", "~"}), "crower_exec=crower search --folder " + shell.Quote("logs") + " " + shell.Quote("~")},
		{"both", notes.SearchBothNote([]string{"needle", "."}), "crower_exec=crower search " + shell.Quote("needle") + " " + shell.Quote(".")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.note, tt.want) {
				t.Errorf("note = %q, want it to contain %q", tt.note, tt.want)
			}
		})
	}
}
