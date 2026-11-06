//go:build fyne

package fyneui

import (
	"image/color"
	"slices"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

func historyTestUI(t *testing.T) (*ui, int) {
	t.Helper()

	var history timeline.Log
	for i := 0; i < 3; i++ {
		history.Add(operation.Create, "greet", "note")
	}

	model := desktop.NewModel(
		catalog.New(nil, []catalog.Command{{Name: "greet", Exec: "echo hi"}}),
		history,
		stubExecutor{},
	)

	front, _ := newTestUI(t, model, stubEditor{}, stubReverter{})
	front.historyTab()

	return front, len(history.AllData)
}

// The history table's cells have to be reachable by the canvas walk for the
// same reason as the command rows: the cell type embeds a container (#406).
func TestHistoryCellsArePaintable(t *testing.T) {
	var history timeline.Log
	for i := 0; i < 3; i++ {
		history.Add(operation.Create, "greet", "note")
	}

	model := desktop.NewModel(
		catalog.New(nil, []catalog.Command{{Name: "greet", Exec: "echo hi"}}),
		history,
		stubExecutor{},
	)

	front, window := newTestUI(t, model, stubEditor{}, stubReverter{})
	window.SetContent(front.content())
	window.Resize(fyneSize())

	front.tabs.SelectIndex(1)
	window.Resize(fyneSize())

	if labels := labelsIn(window.Canvas().Content()); !slices.Contains(labels, "greet") {
		t.Fatalf("the history cell text is not reachable in the rendered tree; labels = %q", labels)
	}
}

// Selecting any cell selects the row it belongs to, not the cell, and opens the
// panel below the table.
func TestHistorySelectsTheWholeRow(t *testing.T) {
	front, rows := historyTestUI(t)

	if rows == 0 {
		t.Fatal("the fixture has no history")
	}

	front.table.OnSelected(widget.TableCellID{Row: 1, Col: 3})

	if front.historyRow != 1 {
		t.Errorf("historyRow = %d, want the row of the clicked cell", front.historyRow)
	}

	if !front.historyBody.Visible() {
		t.Error("selecting a row did not open the detail panel")
	}

	if !strings.Contains(front.historyDetail.Text, "Version") {
		t.Errorf("detail = %q, want the entry's version", front.historyDetail.Text)
	}

	if front.revert.Disabled() {
		t.Error("restore stayed disabled after a row was selected")
	}
}

// Selecting a row must give the history panel a height, not only make it
// visible: Fyne's layouts skip an invisible child, so a panel hidden at build
// time has no size until the tab lays out again.
func TestSelectingAHistoryRowSizesTheDetailPanel(t *testing.T) {
	var history timeline.Log
	for i := 0; i < 3; i++ {
		history.Add(operation.Create, "greet", "note")
	}

	model := desktop.NewModel(
		catalog.New(nil, []catalog.Command{{Name: "greet", Exec: "echo hi"}}),
		history,
		stubExecutor{},
	)

	front, window := newTestUI(t, model, stubEditor{}, stubReverter{})
	window.SetContent(front.content())
	window.Resize(fyneSize())

	if size := front.historyBody.Size(); size.Height != 0 {
		t.Fatalf("the panel already has a size before a selection: %v", size)
	}

	front.table.OnSelected(widget.TableCellID{Row: 0, Col: 1})

	if !front.historyBody.Visible() {
		t.Fatal("selecting a row did not open the panel")
	}

	if size := front.historyBody.Size(); size.Height == 0 {
		t.Fatalf("panel size = %v, want the tab to lay it out", size)
	}
}

// Unselecting closes the panel again, so the table keeps the height.
func TestHistoryUnselectClosesTheDetail(t *testing.T) {
	front, _ := historyTestUI(t)

	front.table.OnSelected(widget.TableCellID{Row: 0, Col: 1})
	front.table.OnUnselected(widget.TableCellID{Row: 0, Col: 1})

	if front.historyRow != -1 {
		t.Errorf("historyRow = %d, want -1", front.historyRow)
	}

	if front.historyBody.Visible() {
		t.Error("the detail panel stayed open after the selection was cleared")
	}

	if front.historyDetail.Text != "" {
		t.Errorf("detail = %q, want it cleared", front.historyDetail.Text)
	}

	if len(front.historyFields.Objects) != 0 {
		t.Error("the field grid stayed populated after the selection was cleared")
	}

	if !front.runHistory.Disabled() || !front.editHistory.Disabled() {
		t.Error("an action stayed enabled without a selected row")
	}

	if !front.revert.Disabled() {
		t.Error("restore stayed enabled without a selected row")
	}
}

// The Details toggle collapses the panel without losing the selection.
func TestHistoryDetailToggle(t *testing.T) {
	front, _ := historyTestUI(t)

	front.table.OnSelected(widget.TableCellID{Row: 0, Col: 0})
	front.setHistoryDetailVisible(false)

	if front.historyBody.Visible() {
		t.Error("the toggle did not collapse the panel")
	}

	if front.historyRow != 0 {
		t.Error("collapsing the panel dropped the selected row")
	}
}

// A cell's background is the row highlight: tinted while its row is selected,
// transparent otherwise.
func TestHistoryCellPaintsTheRowSelection(t *testing.T) {
	cell := newHistoryCell()
	cell.label.SetText("Version 1")

	cell.setSelected(true)

	if cell.bg.FillColor == color.Transparent {
		t.Error("the selected row is not tinted")
	}

	cell.setSelected(false)

	if cell.bg.FillColor != color.Transparent {
		t.Errorf("FillColor = %v, want transparent for an unselected row", cell.bg.FillColor)
	}
}

// historyEntryUI lays out the window over one command and one recorded entry,
// with History selected and that row chosen.
func historyEntryUI(t *testing.T, exec desktop.Executor, op operation.Operation, command, note string) (*ui, fyne.Window) {
	t.Helper()

	var history timeline.Log
	history.Add(op, command, note)

	model := desktop.NewModel(
		catalog.New(nil, []catalog.Command{{Name: "greet", Exec: "echo hi"}}),
		history,
		exec,
	)

	front, window := newTestUI(t, model, stubEditor{}, stubReverter{})
	window.SetContent(front.content())
	window.Resize(fyneSize())
	front.tabs.SelectIndex(1)
	front.table.OnSelected(widget.TableCellID{Row: 0, Col: 1})

	return front, window
}

// The note is shown as labelled fields rather than one key=value line, and the
// original line stays behind the Raw note toggle.
func TestHistoryDetailLabelsTheNoteFields(t *testing.T) {
	command := catalog.Command{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"}
	note := notes.ExecuteCommandNote(command, []string{"a b"})

	front, _ := historyEntryUI(t, stubExecutor{}, operation.Execute, "greet", note)

	labels := labelsIn(front.historyFields)
	for _, want := range []string{"Operation", "Name", "Aliases", "Exec", "Params", "execute", "greet", "hi", "echo hi", "a b"} {
		if !slices.Contains(labels, want) {
			t.Errorf("fields %q do not contain %q", labels, want)
		}
	}

	if front.historyRaw.Visible() {
		t.Error("the raw note is open before the toggle was pressed")
	}

	if !strings.Contains(front.historyRaw.Text, "command_exec=echo hi") {
		t.Errorf("raw note = %q, want the key=value line", front.historyRaw.Text)
	}
}

// Run replays the entry's command with the arguments it recorded.
func TestHistoryRunReplaysTheRecordedParams(t *testing.T) {
	exec := &recordingExecutor{}
	note := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, []string{"a b", "-x"})

	front, _ := historyEntryUI(t, exec, operation.Execute, "greet", note)

	if front.runHistory.Disabled() {
		t.Fatal("run is disabled for an execute entry whose command exists")
	}

	front.runHistoryEntry()

	if !slices.Equal(exec.keys, []string{"greet"}) {
		t.Fatalf("executed %v, want [greet]", exec.keys)
	}

	if !slices.Equal(exec.args[0], []string{"a b", "-x"}) {
		t.Errorf("params = %v, want the recorded arguments", exec.args[0])
	}
}

// Edit clears the Commands filters, selects the entry's command and opens its
// dialog, so the user does not have to find it by hand.
func TestHistoryEditOpensTheCommand(t *testing.T) {
	note := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, nil)

	front, window := historyEntryUI(t, stubExecutor{}, operation.Execute, "greet", note)

	front.model.Filter = "zzz"
	front.filterEntry.SetText("zzz")
	front.list.Refresh()

	front.editHistoryEntry()

	if front.tabs.SelectedIndex() != 0 {
		t.Errorf("selected tab = %d, want Commands", front.tabs.SelectedIndex())
	}

	if front.model.Filter != "" {
		t.Errorf("filter = %q, want it cleared", front.model.Filter)
	}

	if got := front.model.Selected(); got != "greet" {
		t.Errorf("selected command = %q, want greet", got)
	}

	if len(window.Canvas().Overlays().List()) == 0 {
		t.Error("edit did not open a dialog")
	}
}

// An action is offered only where it can work: Run needs an execute entry whose
// command exists, Edit only needs the command to exist.
func TestHistoryActionsFollowTheEntry(t *testing.T) {
	tests := []struct {
		name    string
		op      operation.Operation
		command string
		note    string
		run     bool
		edit    bool
	}{
		{
			"execute with a live command",
			operation.Execute, "greet",
			"main_operation=execute;sub_operation=command;command_name=greet",
			true, true,
		},
		{
			"copy naming a live command",
			operation.Copy, "greet",
			"main_operation=copy",
			false, true,
		},
		{
			"a command that no longer exists",
			operation.Execute, "gone",
			"main_operation=execute;sub_operation=command;command_name=gone",
			false, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			front, _ := historyEntryUI(t, stubExecutor{}, tt.op, tt.command, tt.note)

			if got := !front.runHistory.Disabled(); got != tt.run {
				t.Errorf("run enabled = %v, want %v", got, tt.run)
			}

			if got := !front.editHistory.Disabled(); got != tt.edit {
				t.Errorf("edit enabled = %v, want %v", got, tt.edit)
			}
		})
	}
}

// The Raw note toggle opens and closes the stored line.
func TestHistoryRawNoteToggle(t *testing.T) {
	note := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, nil)

	front, _ := historyEntryUI(t, stubExecutor{}, operation.Execute, "greet", note)

	if front.historyRaw.Visible() {
		t.Fatal("the raw note starts open")
	}

	front.historyRawToggle.OnTapped()
	if !front.historyRaw.Visible() || !front.historyRawShown {
		t.Error("the toggle did not open the raw note")
	}

	front.historyRawToggle.OnTapped()
	if front.historyRaw.Visible() || front.historyRawShown {
		t.Error("the toggle did not close the raw note")
	}
}

// The stored list syntax is not what a reader should see.
func TestHistoryFieldValueFormatsLists(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		want  string
	}{
		{"aliases", "command_alias", "[a b]", "a b"},
		{"steps", "command_steps", `["echo a","echo b"]`, "echo a\necho b"},
		{"params", "command_params", `["a b","-x"]`, "a b\n-x"},
		{"corrupt json is shown as stored", "command_params", "[bad", "[bad"},
		{"a plain value is unchanged", "command_exec", "echo hi", "echo hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := historyFieldValue(tt.key, tt.value); got != tt.want {
				t.Errorf("historyFieldValue(%q, %q) = %q, want %q", tt.key, tt.value, got, tt.want)
			}
		})
	}
}

// Unknown keys are humanised rather than dropped, so a newer note still reads.
func TestHistoryFieldLabel(t *testing.T) {
	tests := map[string]string{
		"main_operation": "Operation",
		"command_name":   "Name",
		"future_field":   "Future Field",
	}

	for key, want := range tests {
		if got := historyFieldLabel(key); got != want {
			t.Errorf("historyFieldLabel(%q) = %q, want %q", key, got, want)
		}
	}
}
