//go:build fyne

package fyneui

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
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
		u.clearHistoryDetail()
	}

	u.runHistory = widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), u.runHistoryEntry)
	u.runHistory.Disable()

	u.editHistory = widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), u.editHistoryEntry)
	u.editHistory.Disable()

	u.revert = widget.NewButtonWithIcon("Restore this state…", theme.HistoryIcon(), u.showRevert)
	u.revert.Disable()

	u.historyHint = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	u.historyHint.Wrapping = fyne.TextWrapWord
	u.historyHint.Hide()

	// The note is shown as labelled fields; the original key=value line stays
	// one click away because it is what a bug report and `list --history` show.
	u.historyFields = container.NewGridWithColumns(2)

	u.historyRaw = widget.NewLabel("")
	u.historyRaw.Wrapping = fyne.TextWrapWord
	u.historyRaw.Hide()

	u.historyRawToggle = widget.NewButton("Raw note", func() {
		u.setHistoryRawVisible(!u.historyRawShown)
	})
	u.historyRawToggle.Importance = widget.LowImportance

	// The detail panel sits below the table, as on the Commands tab, and stays
	// collapsed until an entry is chosen.
	u.historyBody = container.NewVBox(
		u.historyDetail,
		u.historyFields,
		u.historyRawToggle,
		u.historyRaw,
		container.NewHBox(u.runHistory, u.editHistory, u.revert),
		u.historyHint,
	)
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
		if row.Command == "" {
			return "-"
		}

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

	if u.historyFields != nil {
		u.historyFields.RemoveAll()
		u.historyFields.Refresh()
	}

	if u.historyRaw != nil {
		u.historyRaw.SetText("")
	}

	u.setHistoryRawVisible(false)

	if u.runHistory != nil {
		u.runHistory.Disable()
	}

	if u.editHistory != nil {
		u.editHistory.Disable()
	}

	if u.revert != nil {
		u.revert.Disable()
	}

	if u.historyHint != nil {
		u.historyHint.SetText("")
		u.historyHint.Hide()
	}

	u.setHistoryDetailVisible(false)
}

// setHistoryRawVisible is the Raw note toggle.
//
// Like the Details panel it has to refresh the tab root: a label hidden at build
// time was skipped by the layout, so showing it alone gives it no height.
func (u *ui) setHistoryRawVisible(visible bool) {
	u.historyRawShown = visible

	if u.historyRaw == nil {
		return
	}

	if visible {
		u.historyRaw.Show()
	} else {
		u.historyRaw.Hide()
	}

	if u.historyRoot != nil {
		u.historyRoot.Refresh()
	}
}

// selectedHistoryRow is the chosen entry, or false when nothing is selected.
func (u *ui) selectedHistoryRow() (desktop.HistoryRow, bool) {
	if u.historyRow < 0 || u.historyRow >= len(u.historyRows) {
		return desktop.HistoryRow{}, false
	}

	return u.historyRows[u.historyRow], true
}

func (u *ui) showHistoryDetail() {
	row, ok := u.selectedHistoryRow()
	if !ok {
		return
	}

	headline := fmt.Sprintf("Version %d — %s", row.Version, row.Operation)
	if row.Command != "" {
		headline += " · " + row.Command
	}
	headline += " · " + row.When

	u.historyDetail.SetText(headline)
	u.showHistoryFields(row.Fields)
	u.historyRaw.SetText(row.Note)
	u.setHistoryRawVisible(false)
	u.updateHistoryActions(row)
	u.revert.Enable()

	u.setHistoryDetailVisible(true)
}

// showHistoryFields lays the note out as label/value pairs. The keys are the
// stored field names, mapped to words a reader recognises; an unknown key is
// still shown (humanised), so a newer note format does not lose information.
func (u *ui) showHistoryFields(fields []desktop.HistoryField) {
	if u.historyFields == nil {
		return
	}

	u.historyFields.RemoveAll()

	if len(fields) == 0 {
		u.historyFields.Add(widget.NewLabel(""))
		u.historyFields.Add(widget.NewLabel("(no note)"))
		u.historyFields.Refresh()
		return
	}

	for _, field := range fields {
		key := widget.NewLabelWithStyle(historyFieldLabel(field.Key), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		key.Wrapping = fyne.TextWrapWord

		value := widget.NewLabel(historyFieldValue(field.Key, field.Value))
		value.Wrapping = fyne.TextWrapWord

		u.historyFields.Add(key)
		u.historyFields.Add(value)
	}

	u.historyFields.Refresh()
}

// historyFieldLabel is what the grid calls a stored key.
func historyFieldLabel(key string) string {
	if label, ok := historyFieldLabels[key]; ok {
		return label
	}

	words := strings.Split(key, "_")
	for i, word := range words {
		if word == "" {
			continue
		}

		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		words[i] = string(runes)
	}

	return strings.Join(words, " ")
}

var historyFieldLabels = map[string]string{
	"main_operation":    "Operation",
	"sub_operation":     "Kind",
	"command_name":      "Name",
	"command_alias":     "Aliases",
	"command_exec":      "Exec",
	"command_steps":     "Steps",
	"command_shell":     "Shell",
	"command_env":       "Env keys",
	"command_timeout":   "Timeout",
	"command_params":    "Params",
	"old_command_name":  "Old name",
	"old_command_alias": "Old aliases",
	"old_command_exec":  "Old exec",
	"crower_exec":       "Recorded as",
}

// historyFieldValue renders a stored value for a reader: lists are lists, not
// Go's bracket syntax or raw JSON.
func historyFieldValue(key, value string) string {
	switch key {
	case "command_alias", "old_command_alias":
		return strings.TrimSpace(strings.Trim(value, "[]"))
	case "command_steps", "command_params":
		var items []string
		if err := json.Unmarshal([]byte(value), &items); err == nil {
			return strings.Join(items, "\n")
		}
	}

	return value
}

// updateHistoryActions enables only what this entry can do: Run needs an execute
// entry whose command still exists, Edit needs the command to exist. A disabled
// action says why, so a grey button is not a dead end.
func (u *ui) updateHistoryActions(row desktop.HistoryRow) {
	if u.runHistory == nil || u.editHistory == nil {
		return
	}

	command, exists := u.model.Commands.Get(row.Command)

	if exists && row.Operation == string(operation.Execute) {
		u.runHistory.Enable()
	} else {
		u.runHistory.Disable()
	}

	if exists {
		u.editHistory.Enable()
	} else {
		u.editHistory.Disable()
	}

	if u.historyHint != nil {
		u.historyHint.SetText(historyActionHint(row, exists, command))
		if u.historyHint.Text == "" {
			u.historyHint.Hide()
		} else {
			u.historyHint.Show()
		}
	}
}

// historyActionHint explains a disabled action, or says nothing when both apply.
func historyActionHint(row desktop.HistoryRow, exists bool, command catalog.Command) string {
	if !exists {
		return "This command is not in the current setup, so it cannot be run or edited."
	}

	if row.Operation != string(operation.Execute) {
		return "Run replays execute entries; this entry recorded a " + row.Operation + "."
	}

	if !row.ParamsRecorded && len(command.Params) > 0 {
		return "This entry did not record its parameters; Run asks for them."
	}

	return ""
}

// runHistoryEntry runs the managed command again with the parameters the entry
// recorded. It is the safe half of "replay": it never re-creates or deletes.
//
// An entry written before command_params existed has no record of its arguments.
// When the command declares parameters, the user is asked for them rather than
// having the command run with literal placeholders.
func (u *ui) runHistoryEntry() {
	row, ok := u.selectedHistoryRow()
	if !ok {
		return
	}

	command, found := u.model.Commands.Get(row.Command)
	if !found {
		return
	}

	if !row.ParamsRecorded && len(command.Params) > 0 {
		u.askParams(row.Command, command)
		return
	}

	u.model.RunWithParams(row.Command, row.Params)
	u.showOutput()
}

// editHistoryEntry takes the entry's command to the Commands tab and opens its
// edit dialog, so the user does not have to find it by hand.
func (u *ui) editHistoryEntry() {
	row, ok := u.selectedHistoryRow()
	if !ok {
		return
	}

	if _, found := u.model.Commands.Get(row.Command); !found {
		return
	}

	u.selectCommand(row.Command)
	u.showEditFor(row.Command)
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
