package output_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/output"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

func historyLog() timeline.Log {
	var log timeline.Log
	log.Add(operation.Create, "first-command", "first")
	log.Add(operation.Delete, "second-command", "second")

	return log
}

// Every history row used to end with a space, because the last column was
// padded to a fixed width and then followed by a literal space. A table is data
// a caller may compare or pipe, so a trailing space is noise that breaks the
// comparison.
func TestHistoryTableHasNoTrailingSpace(t *testing.T) {
	log := historyLog()

	table := output.HistoryTable(log)

	for _, line := range strings.Split(strings.TrimRight(table, "\n"), "\n") {
		if strings.HasSuffix(line, " ") {
			t.Errorf("line ends with a space: %q", line)
		}
	}

	if !strings.Contains(table, "Timestamp") {
		t.Errorf("table = %q, want the header to say Timestamp", table)
	}
}

// The newest entry is listed first, so `crower list --history` reads the way a
// reader expects.
func TestHistoryTablePutsTheNewestEntryFirst(t *testing.T) {
	table := output.HistoryTable(historyLog())

	if strings.Index(table, "second-command") > strings.Index(table, "first-command") {
		t.Errorf("table = %q, want the newest entry first", table)
	}
}

func TestHistoryTableShowsTheCommandAndADash(t *testing.T) {
	var log timeline.Log
	log.Add(operation.Create, "greet", "a plain note")
	log.Add(operation.Copy, "", "another note")

	table := output.HistoryTable(log)

	for _, want := range []string{"Command", "greet", " - "} {
		if !strings.Contains(table, want) {
			t.Errorf("table = %q, want %q", table, want)
		}
	}
}

func TestHistorySlices(t *testing.T) {
	var log timeline.Log
	for _, name := range []string{"alpha", "beta", "gamma"} {
		log.Add(operation.Create, name, "note")
	}

	first := output.HistoryFirst(log, 2)
	if !strings.Contains(first, "alpha") || !strings.Contains(first, "beta") || strings.Contains(first, "gamma") {
		t.Errorf("first = %q, want the oldest two", first)
	}

	last := output.HistoryLast(log, 2)
	if !strings.Contains(last, "gamma") || !strings.Contains(last, "beta") || strings.Contains(last, "alpha") {
		t.Errorf("last = %q, want the newest two", last)
	}

	if lines := strings.Split(strings.TrimRight(output.HistoryTable(log), "\n"), "\n"); len(lines) != 6 {
		t.Errorf("table has %d lines, want header + border + 3 rows", len(lines))
	}
}
