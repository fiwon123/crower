package output_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/output"
	"github.com/fiwon123/crower/internal/domain/catalog"
)

// The table is both `list`'s result and every prompt's question, so its bytes
// are a contract: row numbers follow the saved order and aliases are joined
// with a comma.
func TestCommandsTable(t *testing.T) {
	commands := catalog.New(
		[]string{"zebra", "alpha"},
		[]catalog.Command{
			{Name: "alpha", Aliases: []string{"a", "b"}, Exec: "echo a"},
			{Name: "zebra", Exec: "echo z"},
		},
	)

	got := output.CommandsTable(commands)

	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")

	if len(lines) != 5 {
		t.Fatalf("table has %d lines, want 5: %q", len(lines), got)
	}

	if !strings.HasPrefix(lines[1], "Row Name") {
		t.Errorf("header = %q, want the column titles", lines[1])
	}

	if !strings.HasPrefix(lines[3], "0   zebra") || !strings.HasSuffix(lines[3], "echo z") {
		t.Errorf("first row = %q, want zebra at row 0", lines[3])
	}

	if !strings.Contains(lines[4], "a,b") {
		t.Errorf("second row = %q, want the aliases joined", lines[4])
	}
}

func TestCommandsTableWithNoCommands(t *testing.T) {
	got := output.CommandsTable(catalog.New(nil, nil))

	if lines := strings.Split(strings.TrimRight(got, "\n"), "\n"); len(lines) != 3 {
		t.Errorf("table = %q, want only the header", got)
	}
}
