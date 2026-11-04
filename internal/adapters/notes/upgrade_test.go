package notes_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
)

func TestUpgradeNote(t *testing.T) {
	got := notes.UpgradeNote()

	for _, want := range []string{"main_operation=upgrade", "crower_exec=crower --upgrade"} {
		if !strings.Contains(got, want) {
			t.Errorf("note = %q, want it to contain %q", got, want)
		}
	}
}
