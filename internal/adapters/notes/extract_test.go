package notes_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/shell"
)

func TestExtractNoteRecordReplayableCommands(t *testing.T) {
	plain := notes.ExtractNote("", []string{"a.zip"})
	if want := "crower_exec=crower extract " + shell.Quote("a.zip"); !strings.Contains(plain, want) {
		t.Errorf("note = %q, want it to contain %q", plain, want)
	}

	withOut := notes.ExtractNote("./build", []string{"a.zip"})
	want := "crower_exec=crower extract --out " + shell.Quote("./build") + " " + shell.Quote("a.zip")
	if !strings.Contains(withOut, want) {
		t.Errorf("note = %q, want it to contain %q", withOut, want)
	}

	if !strings.HasPrefix(withOut, "main_operation=extract") {
		t.Errorf("note = %q, want the extract operation", withOut)
	}
}
