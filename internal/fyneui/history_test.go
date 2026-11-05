//go:build fyne

package fyneui

import (
	"image/color"
	"slices"
	"strings"
	"testing"

	"fyne.io/fyne/v2/widget"

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
