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

// HistorySummary is one entry as a collapsed card shows it.
//
// It carries no note: a page is built from summaries without parsing the log, so
// a long history costs the same to list as a short one.
type HistorySummary struct {
	Version   int
	When      string
	Operation string
	Command   string
}

// HistoryDetail is one entry's expanded view: the note split into labelled
// fields, the parameters an execute recorded, and the original line. Position is
// the entry's distance from the newest (0), so Restore knows how many newer
// entries it drops.
type HistoryDetail struct {
	Version        int
	Position       int
	When           string
	Operation      string
	Command        string
	Note           string
	Fields         []HistoryField
	Params         []string
	ParamsRecorded bool
}

// HistoryText renders the history screen for another front end. The Fyne
// window shows expandable cards built from summaries, but the text rendering
// stays here because the terminal front end and its test use it.
func (m *Model) HistoryText() string {
	return m.historyText()
}

// HistoryPage returns up to limit summaries starting offset entries back from
// the newest: offset 0 is the newest operation.
//
// Only the requested page is formatted, which is what keeps a long log from
// doing work the window is not showing.
func (m *Model) HistoryPage(offset, limit int) []HistorySummary {
	total := len(m.History.AllData)
	if offset < 0 {
		offset = 0
	}

	if limit < 0 {
		limit = 0
	}

	if offset >= total || limit == 0 {
		return []HistorySummary{}
	}

	end := min(total, offset+limit)
	summaries := make([]HistorySummary, 0, end-offset)

	for position := offset; position < end; position++ {
		entry := m.History.AllData[total-1-position]

		summaries = append(summaries, HistorySummary{
			Version:   entry.Version,
			When:      historyWhen(entry.Timestemp),
			Operation: string(entry.Operation),
			Command:   entry.CommandName,
		})
	}

	return summaries
}

// HistoryDetail returns the expanded view of one entry by version. A card holds
// the version, so an entry appended after the page was built cannot shift it.
func (m *Model) HistoryDetail(version int) (HistoryDetail, bool) {
	for position := len(m.History.AllData) - 1; position >= 0; position-- {
		entry := m.History.AllData[position]
		if entry.Version != version {
			continue
		}

		fields := notes.Fields(entry.Note)
		historyFields := make([]HistoryField, 0, len(fields))

		for _, field := range fields {
			historyFields = append(historyFields, HistoryField{Key: field.Key, Value: field.Value})
		}

		params, paramsRecorded := notes.CommandParams(entry.Note)

		return HistoryDetail{
			Version:        entry.Version,
			Position:       len(m.History.AllData) - 1 - position,
			When:           historyWhen(entry.Timestemp),
			Operation:      string(entry.Operation),
			Command:        entry.CommandName,
			Note:           notes.Display(entry.Note),
			Fields:         historyFields,
			Params:         params,
			ParamsRecorded: paramsRecorded,
		}, true
	}

	return HistoryDetail{}, false
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
