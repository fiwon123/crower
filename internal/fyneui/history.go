//go:build fyne

package fyneui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
)

// historyHeaders is the table's column order. Version is last because the
// number is an identifier, not something a reader scans for.
var historyHeaders = []string{"When", "Operation", "Command", "Version"}

// historyTab renders every recorded operation as a table, newest first, with a
// restore action.
//
// It exists because the old History dialog showed 21 text lines with no way to
// act on them: the log is what makes crower undoable, so it gets a view.
func (u *ui) historyTab() fyne.CanvasObject {
	u.historyRows = u.model.HistoryRows()

	u.historyDetail = widget.NewLabel("Select an entry to see what it recorded.")
	u.historyDetail.Wrapping = fyne.TextWrapWord

	u.table = widget.NewTableWithHeaders(
		func() (int, int) { return len(u.historyRows), len(historyHeaders) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label, ok := cell.(*widget.Label)
			if !ok || id.Row < 0 || id.Row >= len(u.historyRows) {
				return
			}

			label.SetText(historyCell(u.historyRows[id.Row], id.Col))
		},
	)

	u.table.ShowHeaderRow = true
	u.table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	}
	u.table.UpdateHeader = func(id widget.TableCellID, cell fyne.CanvasObject) {
		label, ok := cell.(*widget.Label)
		if !ok || id.Col < 0 || id.Col >= len(historyHeaders) {
			return
		}

		label.SetText(historyHeaders[id.Col])
	}

	u.table.SetColumnWidth(0, 160)
	u.table.SetColumnWidth(1, 140)
	u.table.SetColumnWidth(2, 220)
	u.table.SetColumnWidth(3, 90)

	u.table.OnSelected = func(id widget.TableCellID) {
		u.historyRow = id.Row
		u.showHistoryDetail()
	}
	u.table.OnUnselected = func(widget.TableCellID) {
		u.historyRow = -1
		u.revert.Disable()
		u.historyDetail.SetText("Select an entry to see what it recorded.")
	}

	u.revert = widget.NewButtonWithIcon("Restore this state…", theme.HistoryIcon(), u.showRevert)
	u.revert.Disable()

	header := widget.NewLabel("Every recorded operation, newest first. Restoring an entry puts the command set back to what that entry recorded and drops the newer ones.")
	header.Wrapping = fyne.TextWrapWord

	return container.NewBorder(
		container.NewVBox(header, u.historyDetail),
		container.NewHBox(u.revert),
		nil,
		nil,
		u.table,
	)
}

func historyCell(row desktop.HistoryRow, column int) string {
	switch column {
	case 0:
		return row.When
	case 1:
		return row.Operation
	case 2:
		return row.Command
	case 3:
		return strconv.Itoa(row.Version)
	default:
		return ""
	}
}

// refreshHistory rebuilds the table's snapshot. A command edit changes the log,
// and refreshing a hidden table on every keystroke would be wasted work.
func (u *ui) refreshHistory() {
	u.historyRows = u.model.HistoryRows()
	u.historyRow = -1

	if u.table != nil {
		u.table.Refresh()
		u.table.UnselectAll()
	}

	if u.revert != nil {
		u.revert.Disable()
	}

	if u.historyDetail != nil {
		u.historyDetail.SetText("Select an entry to see what it recorded.")
	}
}

func (u *ui) showHistoryDetail() {
	if u.historyRow < 0 || u.historyRow >= len(u.historyRows) {
		return
	}

	row := u.historyRows[u.historyRow]

	note := row.Note
	if note == "" {
		note = "(no note)"
	}

	u.historyDetail.SetText(fmt.Sprintf("Version %d — %s, %s\n%s", row.Version, row.Operation, row.Command, note))
	u.revert.Enable()
}

// showRevert asks before restoring, then applies the change through the
// Reverter port.
func (u *ui) showRevert() {
	if u.historyRow < 0 || u.historyRow >= len(u.historyRows) {
		return
	}

	row := u.historyRows[u.historyRow]

	// RevertTo(steps) restores the entry steps back from the newest, so row 0
	// (the newest) is one step and the entry itself stays in the log.
	steps := u.historyRow + 1

	dropped := steps - 1
	detail := "No newer entries are dropped."
	if dropped == 1 {
		detail = "One newer entry is dropped."
	} else if dropped > 1 {
		detail = fmt.Sprintf("%d newer entries are dropped.", dropped)
	}

	question := fmt.Sprintf(
		"Restore the command set recorded at version %d (%s, %s)?\n\n%s",
		row.Version, row.Operation, row.Command, detail,
	)

	dialog.ShowConfirm("Restore command set", question, func(ok bool) {
		if !ok {
			return
		}

		u.applyRevert(steps)
	}, u.window)
}

// applyRevert is the confirmation's submit path.
func (u *ui) applyRevert(steps int) bool {
	commands, history, err := u.reverter.Revert(steps)
	if err != nil {
		showError(err.Error(), u.window)
		return false
	}

	u.reload(commands, history)
	return true
}
