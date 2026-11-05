//go:build fyne

package fyneui

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

	u.historyDetail = widget.NewLabel("")
	u.historyDetail.Wrapping = fyne.TextWrapWord

	u.table = widget.NewTableWithHeaders(
		func() (int, int) { return len(u.historyRows), len(historyHeaders) },
		func() fyne.CanvasObject { return newHistoryCell() },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			view, ok := cell.(*historyCellView)
			if !ok || id.Row < 0 || id.Row >= len(u.historyRows) {
				return
			}

			view.label.SetText(historyCell(u.historyRows[id.Row], id.Col))
			// The selection is a row, so every cell of the selected row is
			// painted the same; a cell on its own means nothing here.
			view.setSelected(id.Row == u.historyRow)
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
		u.table.Refresh()
		u.showHistoryDetail()
	}
	u.table.OnUnselected = func(widget.TableCellID) {
		u.historyRow = -1
		u.table.Refresh()
		u.revert.Disable()
		u.clearHistoryDetail()
	}

	u.revert = widget.NewButtonWithIcon("Restore this state…", theme.HistoryIcon(), u.showRevert)
	u.revert.Disable()

	// The detail panel sits below the table, as on the Commands tab, and stays
	// collapsed until an entry is chosen.
	u.historyBody = container.NewVBox(u.historyDetail, container.NewHBox(u.revert))
	u.historyBody.Hide()

	u.historyToggle = widget.NewButtonWithIcon("Details", theme.MenuDropDownIcon(), func() {
		u.setHistoryDetailVisible(!u.historyDetailShown)
	})
	u.historyToggle.Importance = widget.LowImportance

	header := widget.NewLabel("Every recorded operation, newest first. Restoring an entry puts the command set back to what that entry recorded and drops the newer ones.")
	header.Wrapping = fyne.TextWrapWord

	u.historyRoot = container.NewBorder(
		header,
		container.NewVBox(u.historyToggle, u.historyBody),
		nil,
		nil,
		u.table,
	)

	return u.historyRoot
}

// historyCellView is one table cell: a label over a background the row
// selection tints. widget.Table recycles cells, so UpdateCell re-points the
// label and the tint on every pass.
//
// Like commandRow, it is a widget holding a container: a struct that embeds a
// container is opaque to the canvas walk and paints nothing (#406).
type historyCellView struct {
	widget.BaseWidget

	bg    *canvas.Rectangle
	label *widget.Label
	body  *fyne.Container
}

func newHistoryCell() *historyCellView {
	cell := &historyCellView{
		bg:    canvas.NewRectangle(color.Transparent),
		label: widget.NewLabel(""),
	}

	cell.body = container.NewMax(cell.bg, cell.label)
	cell.ExtendBaseWidget(cell)

	return cell
}

func (c *historyCellView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.body)
}

func (c *historyCellView) setSelected(selected bool) {
	if selected {
		c.bg.FillColor = theme.Color(theme.ColorNameSelection)
	} else {
		c.bg.FillColor = color.Transparent
	}

	c.bg.Refresh()
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

	u.clearHistoryDetail()
}

// setHistoryDetailVisible is the Details toggle below the table.
//
// As on the Commands tab, showing a panel that was hidden at build time needs the
// tab root to lay out again: Fyne's layouts skip invisible children, so without
// the refresh the panel would keep the zero size it was never given.
func (u *ui) setHistoryDetailVisible(visible bool) {
	u.historyDetailShown = visible

	if u.historyBody == nil {
		return
	}

	if visible {
		u.historyBody.Show()
	} else {
		u.historyBody.Hide()
	}

	if u.historyRoot != nil {
		u.historyRoot.Refresh()
	}
}

// clearHistoryDetail empties and collapses the panel, so a refresh cannot leave
// another entry's text on screen.
func (u *ui) clearHistoryDetail() {
	if u.historyDetail != nil {
		u.historyDetail.SetText("")
	}

	u.setHistoryDetailVisible(false)
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
	u.setHistoryDetailVisible(true)
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
