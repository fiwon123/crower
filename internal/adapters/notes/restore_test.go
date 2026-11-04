package notes_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/shell"
)

func TestRestoreNoteRecordsTheKey(t *testing.T) {
	built := notes.RestoreNote([]string{"greet"})

	want := "crower_exec=crower restore " + shell.Quote("greet")
	if !strings.Contains(built, want) {
		t.Errorf("note = %q, want it to contain %q", built, want)
	}
}

// An interactive restore has no argument; the recorded line is the bare
// invocation, which is what ran.
func TestRestoreNoteWithoutAKey(t *testing.T) {
	built := notes.RestoreNote(nil)

	if want := "crower_exec=crower restore"; !strings.HasSuffix(built, want) {
		t.Errorf("note = %q, want it to end with %q", built, want)
	}
}
