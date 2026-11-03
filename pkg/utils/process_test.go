package utils

import (
	"strings"
	"testing"
)

func TestListProcessesReturnsRows(t *testing.T) {
	rows, err := ListProcesses("", false)
	if err != nil {
		t.Fatalf("ListProcesses: %v", err)
	}

	if len(rows) == 0 {
		t.Fatal("ListProcesses returned no processes at all")
	}

	for _, row := range rows {
		if row.PID <= 0 {
			t.Errorf("process row has no pid: %+v", row)
		}

		if row.Name == "" {
			t.Errorf("process row has no name: %+v", row)
		}
	}
}

func TestListProcessesFiltersByName(t *testing.T) {
	all, err := ListProcesses("", false)
	if err != nil {
		t.Fatalf("ListProcesses: %v", err)
	}

	// Filter by a name the listing actually contains, so the test does not
	// depend on which processes the machine happens to run.
	name := all[0].Name
	rows, err := ListProcesses(strings.ToUpper(name), false)
	if err != nil {
		t.Fatalf("ListProcesses: %v", err)
	}

	if len(rows) == 0 {
		t.Fatalf("filter %q matched nothing, but %q is running", name, name)
	}

	for _, row := range rows {
		if !strings.Contains(strings.ToLower(row.Name), strings.ToLower(name)) {
			t.Errorf("row %q does not contain the filter %q", row.Name, name)
		}
	}
}
