//go:build fyne

package fyneui

import (
	"image/color"
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
