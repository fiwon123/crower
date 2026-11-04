package notes_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/shell"
)

func TestTransferNotesRecordReplayableCommands(t *testing.T) {
	tests := []struct {
		name string
		note string
		want string
		sub  string
	}{
		{
			"copy files",
			notes.CopyNote(operation.File, []string{"a.txt", "b.txt", "dest"}),
			"crower_exec=crower copy " + shell.Quote("a.txt") + " " + shell.Quote("b.txt") + " " + shell.Quote("dest"),
			"sub_operation=file",
		},
		{
			"copy folders",
			notes.CopyNote(operation.Folder, []string{"dir", "dest"}),
			"crower_exec=crower copy " + shell.Quote("dir") + " " + shell.Quote("dest"),
			"sub_operation=folder",
		},
		{
			"move mixed",
			notes.MoveNote(operation.FileAndFolder, []string{"a.txt", "dir", "dest"}),
			"crower_exec=crower move " + shell.Quote("a.txt") + " " + shell.Quote("dir") + " " + shell.Quote("dest"),
			"sub_operation=file_and_folder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.note, tt.want) {
				t.Errorf("note = %q, want it to contain %q", tt.note, tt.want)
			}

			if !strings.Contains(tt.note, tt.sub) {
				t.Errorf("note = %q, want it to contain %q", tt.note, tt.sub)
			}
		})
	}
}
