package output

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// HistoryTable renders the whole log, newest first.
func HistoryTable(log timeline.Log) string {
	var builder strings.Builder

	builder.WriteString(historyHeader())

	line := 0

	for i := len(log.AllData) - 1; i >= 0; i-- {
		entry := log.AllData[i]
		builder.WriteString(historyRow(line, entry))
		line++
	}

	return builder.String()
}

// HistoryFirst renders the oldest entries up to steps.
func HistoryFirst(log timeline.Log, steps int) string {
	var builder strings.Builder

	builder.WriteString(historyHeader())

	line := 0

	for i := steps - 1; i >= 0; i-- {
		builder.WriteString(historyRow(line, log.AllData[i]))
		line++
	}

	return builder.String()
}

// HistoryLast renders the newest entries down to steps.
func HistoryLast(log timeline.Log, steps int) string {
	var builder strings.Builder

	builder.WriteString(historyHeader())

	line := 0

	for i := len(log.AllData) - 1; i >= len(log.AllData)-steps; i-- {
		builder.WriteString(historyRow(line, log.AllData[i]))
		line++
	}

	return builder.String()
}

func historyRow(line int, entry timeline.Content) string {
	command := entry.CommandName
	if command == "" {
		command = "-"
	}

	return fmt.Sprintf("%-4d %-8d %-16s %-32s %-20s %s\n",
		line, entry.Version, entry.File, entry.Timestemp, command, notes.Display(entry.Note))
}

func historyHeader() string {
	var builder strings.Builder

	builder.WriteString("----------------------------------------------------------------------------\n")
	fmt.Fprintf(&builder, "%-4s %-8s %-16s %-32s %-20s %s\n", "line", "Version", "File", "Timestamp", "Command", "Note")
	builder.WriteString("----------------------------------------------------------------------------\n")

	return builder.String()
}
