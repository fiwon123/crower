package extractdom

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/shell"
)

// The --out value is quoted for the recording OS: a hand-written double quote
// broke on a path that contains one.
func TestExtractNoteQuotesTheOutPath(t *testing.T) {
	out := `C:\My Files`
	built := newExtractNote(out, []string{"a.tar"})

	want := "crower_exec=crower extract --out " + shell.Quote(out) + " " + shell.Quote("a.tar")
	if !strings.Contains(built, want) {
		t.Errorf("note = %q, want it to contain %q", built, want)
	}
}
