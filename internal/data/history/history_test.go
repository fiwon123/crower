package history

import (
	"slices"
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
	data.Add(operations.Copy, "", "main_operation=copy;crower_exec=crower copy a b")

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
	data.Add(operations.Create, "first-command", "first")
	data.Add(operations.Delete, "second-command", "second")

	list := data.GetList()

	if strings.Index(list, "second") > strings.Index(list, "first") {
		t.Errorf("GetList() = %q, want the newest entry first", list)
	}
}

// The entry records which command the operation touched. Without it, the
// `--last` family resolved an empty name and could never act.
func TestAddRecordsTheCommandName(t *testing.T) {
	data := New(t.TempDir())
	data.Add(operations.Execute, "greet", "main_operation=execute;command_name=greet")

	last := data.GetLast()
	if last == nil {
		t.Fatal("Add did not append an entry")
	}

	if last.CommandName != "greet" {
		t.Errorf("CommandName = %q, want greet", last.CommandName)
	}
}

// The table names the command an entry touched. The note itself may contain
// that text, so the fixture uses a note without it.
func TestGetListShowsTheCommand(t *testing.T) {
	data := New("history")

	data.Add(operations.Create, "greet", "a plain note")
	data.Add(operations.Copy, "", "another note")

	list := data.GetList()

	if !strings.Contains(list, "Command") {
		t.Errorf("GetList() = %q, want a Command column", list)
	}

	if !strings.Contains(list, "greet") {
		t.Errorf("GetList() = %q, want the command name", list)
	}

	if !strings.Contains(list, "-") {
		t.Errorf("GetList() = %q, want a dash for the command-less entry", list)
	}
}

// One entry is dropped and the rest keep their order.
func TestRemoveContentDropsOneEntry(t *testing.T) {
	data := New(t.TempDir())
	data.Add(operations.Create, "first", "one")
	data.Add(operations.Update, "second", "two")
	data.Add(operations.Delete, "third", "three")

	target := data.AllData[1]

	if ok := data.RemoveContent(target); !ok {
		t.Fatal("RemoveContent of a present entry returned false")
	}

	remaining := []string{}

	for _, entry := range data.AllData {
		remaining = append(remaining, entry.CommandName)
	}

	if want := []string{"first", "third"}; !slices.Equal(remaining, want) {
		t.Errorf("history = %v, want %v", remaining, want)
	}
}

// A stale entry that is not in the history is not a reason to discard the
// history. The call reports the miss and changes nothing (#218).
func TestRemoveContentLeavesTheHistoryOnAMiss(t *testing.T) {
	data := New(t.TempDir())
	data.Add(operations.Create, "first", "one")
	data.Add(operations.Update, "second", "two")

	before := slices.Clone(data.AllData)

	if ok := data.RemoveContent(Content{Version: 99}); ok {
		t.Fatal("RemoveContent of a missing entry returned true")
	}

	if !slices.Equal(data.AllData, before) {
		t.Errorf("history changed on a miss: %v, want %v", data.AllData, before)
	}
}
