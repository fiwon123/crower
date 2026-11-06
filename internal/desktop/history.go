package desktop

import (
	"time"

	"github.com/fiwon123/crower/internal/adapters/notes"
)

// HistoryField is one key/value pair of an entry's note, in the order it was
// written, with the value decoded for display.
type HistoryField struct {
	Key   string
	Value string
}

// HistoryRow is one history entry rendered as a table row.
//
// It is a presentation shape, not the stored entry: the timestamp is formatted,
// the note is shown the way the CLI shows it, and the note is also split into
// its fields so a front end can label them. Params are what an execute entry ran
// with; ParamsRecorded is false for an entry written before the field existed (or
// whose field is corrupt), which is not the same as "ran with none".
type HistoryRow struct {
	Version        int
	When           string
	Operation      string
	Command        string
	Note           string
	Fields         []HistoryField
	Params         []string
	ParamsRecorded bool
}

// HistoryText renders the history screen for another front end. The Fyne
// window shows a table built from HistoryRows, but the text rendering stays
// here because the terminal front end and its test use it.
func (m *Model) HistoryText() string {
	return m.historyText()
}

// HistoryRows returns the log newest first, so row 0 is the most recent
// operation and the state a restore would leave first.
func (m *Model) HistoryRows() []HistoryRow {
	rows := make([]HistoryRow, 0, len(m.History.AllData))

	for i := len(m.History.AllData) - 1; i >= 0; i-- {
		entry := m.History.AllData[i]

		command := entry.CommandName

		fields := notes.Fields(entry.Note)
		historyFields := make([]HistoryField, 0, len(fields))

		for _, field := range fields {
			historyFields = append(historyFields, HistoryField{Key: field.Key, Value: field.Value})
		}

		params, paramsRecorded := notes.CommandParams(entry.Note)

		rows = append(rows, HistoryRow{
			Version:        entry.Version,
			When:           historyWhen(entry.Timestemp),
			Operation:      string(entry.Operation),
			Command:        command,
			Note:           notes.Display(entry.Note),
			Fields:         historyFields,
			Params:         params,
			ParamsRecorded: paramsRecorded,
		})
	}

	return rows
}

// historyWhen renders a stored RFC3339 timestamp as local date and time. A
// value that does not parse is shown as stored rather than dropped: the log
// belongs to the user even when a hand edit broke one field.
func historyWhen(raw string) string {
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}

	return parsed.Local().Format("2006-01-02 15:04")
}
