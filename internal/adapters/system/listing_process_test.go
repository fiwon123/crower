package system

import (
	"errors"
	"testing"

	"github.com/fiwon123/crower/internal/application/listcontent"
	"github.com/fiwon123/crower/pkg/utils"
)

// The adapter maps the process rows the host reports to the result's rows, so a
// JSON key or field rename cannot silently drop a column.
func TestProcessesMapsTheRows(t *testing.T) {
	listings := NewListings()
	listings.listProcesses = func(part string, skipUnknown bool) ([]utils.ProcessInfo, error) {
		if part != "go" {
			t.Errorf("filter = %q, want go", part)
		}

		if !skipUnknown {
			t.Error("skipUnknown = false, want the rows with a known path")
		}

		return []utils.ProcessInfo{
			{PID: 7, Name: "go", Path: "/usr/bin/go"},
		}, nil
	}

	rows, err := listings.Processes("go")
	if err != nil {
		t.Fatalf("Processes: %v", err)
	}

	want := []listcontent.ProcessRow{{PID: 7, Name: "go", Path: "/usr/bin/go"}}
	if len(rows) != len(want) || rows[0] != want[0] {
		t.Errorf("rows = %+v, want %+v", rows, want)
	}
}

func TestProcessesWithNoMatchesReturnsAnEmptySlice(t *testing.T) {
	listings := NewListings()
	listings.listProcesses = func(string, bool) ([]utils.ProcessInfo, error) {
		return nil, nil
	}

	rows, err := listings.Processes("")
	if err != nil {
		t.Fatalf("Processes: %v", err)
	}

	if rows == nil || len(rows) != 0 {
		t.Errorf("rows = %#v, want an empty slice so JSON encodes []", rows)
	}
}

func TestProcessesPassesTheErrorUp(t *testing.T) {
	boom := errors.New("boom")

	listings := NewListings()
	listings.listProcesses = func(string, bool) ([]utils.ProcessInfo, error) {
		return nil, boom
	}

	if _, err := listings.Processes(""); !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}
}
