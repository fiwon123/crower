package operation

import "testing"

// The values are not internal names: they are written to historyhelper.json as
// the Operation field of every entry, and to a history note as
// `main_operation=<value>`. Renaming one makes the history of every existing
// user unreadable, so the wire form is pinned here.
func TestMainOperationValues(t *testing.T) {
	tests := []struct {
		op   Operation
		want string
	}{
		{Execute, "execute"},
		{Create, "create"},
		{Delete, "delete"},
		{Update, "update"},
		{List, "list"},
		{Reset, "reset"},
		{Open, "open"},
		{Revert, "revert"},
		{Search, "search"},
		{Extract, "extract"},
		{Copy, "copy"},
		{Move, "move"},
		{Restore, "restore"},
		{Check, "check"},
		{Upgrade, "upgrade"},
	}

	for _, tt := range tests {
		if got := string(tt.op); got != tt.want {
			t.Errorf("operation = %q, want %q", got, tt.want)
		}
	}
}

func TestSubOperationValues(t *testing.T) {
	tests := []struct {
		op   SubOperation
		want string
	}{
		{Command, "command"},
		{Script, "script"},
		{Process, "process"},
		{File, "file"},
		{Folder, "folder"},
		{FileAndFolder, "file_and_folder"},
		{History, "history"},
		{SysPath, "syspath"},
		{System, "system"},
		{Last, "last"},
		{Browser, "browser"},
	}

	for _, tt := range tests {
		if got := string(tt.op); got != tt.want {
			t.Errorf("sub-operation = %q, want %q", got, tt.want)
		}
	}
}

// Two operations sharing a value would make a history entry ambiguous: the
// reader cannot tell `copy` from `move`, and a `--last` lookup would match
// whichever row came last.
func TestOperationValuesAreUnique(t *testing.T) {
	count := map[string]int{}

	for _, op := range []Operation{
		Execute, Create, Delete, Update, List, Reset, Open,
		Revert, Search, Extract, Copy, Move, Restore, Check, Upgrade,
	} {
		count[string(op)]++
	}

	for value, n := range count {
		if n > 1 {
			t.Errorf("main operation value %q is declared %d times", value, n)
		}
	}
}

func TestOperationValuesAreNotEmpty(t *testing.T) {
	for _, op := range []Operation{
		Execute, Create, Delete, Update, List, Reset, Open,
		Revert, Search, Extract, Copy, Move, Restore, Check, Upgrade,
	} {
		if op == "" {
			t.Errorf("a main operation has an empty value")
		}
	}

	for _, op := range []SubOperation{
		Command, Script, Process, File, Folder, FileAndFolder,
		History, SysPath, System, Last, Browser,
	} {
		if op == "" {
			t.Errorf("a sub-operation has an empty value")
		}
	}
}
