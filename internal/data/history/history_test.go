package history

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/data/operations"
)

// Every history row used to end with a space, because the last column was
// padded to a fixed width and then followed by a literal space. A table is data
// a caller may compare or pipe, so a trailing space is noise that breaks the
// comparison.
func TestGetListHasNoTrailingSpace(t *testing.T) {
	data := New(t.TempDir())
	data.Add(operations.Copy, "main_operation=copy;crower_exec=crower copy a b")

	list := data.GetList()

	if !strings.Contains(list, "crower copy a b") {
		t.Fatalf("GetList() = %q, want it to contain the note", list)
	}

	for _, line := range strings.Split(strings.TrimRight(list, "\n"), "\n") {
		if strings.HasSuffix(line, " ") {
			t.Errorf("line ends with a space: %q", line)
		}
	}

	if !strings.Contains(list, "Timestamp") {
		t.Errorf("GetList() = %q, want the header to say Timestamp", list)
	}
}

// The newest entry is listed first, so `crower list --history` reads the way a
// reader expects.
func TestGetListPutsTheNewestEntryFirst(t *testing.T) {
	data := New(t.TempDir())
	data.Add(operations.Create, "first")
	data.Add(operations.Delete, "second")

	list := data.GetList()

	if strings.Index(list, "second") > strings.Index(list, "first") {
		t.Errorf("GetList() = %q, want the newest entry first", list)
	}
}
