//go:build fyne

package fyneui

import (
	"context"
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
)

// newGrowingEntry builds a multi-line field that shows its text instead of
// hiding it behind a scrollbar: it wraps long lines and grows to the number of
// rows the text needs, between minRows and maxRows.
//
// Past maxRows it keeps its own scroll, which is the one place a wheel event
// still stops at a field — an entry tall enough to hold a script is worth that,
// and the dialogs put the whole form in a scroll so nothing is unreachable.
func newGrowingEntry(minRows, maxRows int) *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapWord
	entry.SetMinRowsVisible(minRows)

	entry.OnChanged = func(text string) {
		entry.SetMinRowsVisible(growingRows(text, minRows, maxRows))
	}

	return entry
}

// growingRows estimates how many rows text needs at the dialog's width, so a
// script shows every line it has up to the cap rather than a fixed window into
// it. The estimate only drives MinRowsVisible: if it is low the field scrolls a
// little, and if it is high a spare empty row shows.
func growingRows(text string, minRows, maxRows int) int {
	const charsPerRow = 72

	rows := 0

	for _, line := range strings.Split(text, "\n") {
		rows += max(1, (len(line)+charsPerRow-1)/charsPerRow)
	}

	return min(max(rows, minRows), maxRows)
}

// newScrollFreeEntry builds a single-line field that does not swallow the mouse
// wheel.
//
// Fyne dispatches a wheel event to the topmost Scrollable under the cursor and
// never bubbles it to a parent (internal/driver/glfw window.processMouseScrolled).
// Every widget.Entry builds an internal scroll unless its wrapping is
// TextWrapOff and its Scroll is container.ScrollNone (widget/entry.go), and the
// constructors' default wrapping is not TextWrapOff — so a plain entry under the
// cursor freezes the scrollbar of the container it sits in. These two settings
// make the entry's renderer drop that internal scroll, and the wheel reaches the
// enclosing scroller.
//
// The trade-off is that a value wider than the field clips instead of scrolling
// inside it. That is what lets the page scroll, which is the reported bug.
func newScrollFreeEntry() *widget.Entry {
	entry := widget.NewEntry()
	entry.Wrapping = fyne.TextWrapOff
	entry.Scroll = container.ScrollNone

	return entry
}

// toolResult is one card's own result area. It stays hidden until the card runs
// something, then opens with the outcome; the toggle hides it again, which is
// what gives the forms back their space.
type toolResult struct {
	panel        *block
	label        *widget.Label
	toggle       *widget.Button
	cancelButton *widget.Button
	body         *fyne.Container
	area         *fyne.Container
	shown        bool
	status       string
	cancel       context.CancelFunc
}

func newToolResult() *toolResult {
	result := &toolResult{panel: newBlock("")}
	result.label = result.panel.label
	result.panel.setAccent(theme.ColorNameSeparator)

	result.cancelButton = widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), result.cancelRun)
	result.cancelButton.Importance = widget.LowImportance
	result.cancelButton.Hide()

	result.body = container.NewVBox(result.panel.object(), result.cancelButton)
	result.body.Hide()

	result.toggle = widget.NewButton("Hide result", func() {
		result.setVisible(!result.shown)
	})
	result.toggle.Importance = widget.LowImportance

	result.area = container.NewVBox(result.toggle, result.body)
	result.area.Hide()

	return result
}

// object is the collapsible area a card places below its actions.
func (r *toolResult) object() fyne.CanvasObject {
	return r.area
}

func (r *toolResult) setVisible(visible bool) {
	r.shown = visible

	if visible {
		r.body.Show()
		r.toggle.SetText("Hide result")
		return
	}

	r.body.Hide()
	r.toggle.SetText("Show result")
}

// report shows a successful outcome in green, and opens the area: a result
// square appears because an action ran, not before.
func (r *toolResult) report(text string) {
	r.reportWith(text, theme.ColorNameSuccess)
}

// reportWith shows an outcome and colours the panel's bar by what happened:
// green for success, amber for a stop, red for a failure.
func (r *toolResult) reportWith(text string, color fyne.ThemeColorName) {
	r.label.SetText(text)
	r.panel.setAccent(color)
	r.area.Show()
	r.setVisible(true)
}

// start opens the area for a run and arms the cancel action. A nil cancel means
// the run cannot be stopped, so no Cancel button is offered.
func (r *toolResult) start(status string, cancel context.CancelFunc) {
	r.status = status
	r.cancel = cancel

	r.label.SetText(status)
	r.panel.setAccent(theme.ColorNameForeground)
	r.area.Show()
	r.setVisible(true)

	if cancel != nil && r.cancelButton != nil {
		r.cancelButton.Show()
		r.cancelButton.Enable()
	}
}

// stop ends a run: the cancel action goes away and the result stays.
func (r *toolResult) stop() {
	r.cancel = nil

	if r.cancelButton != nil {
		r.cancelButton.Hide()
	}
}

// cancelRun asks the running action to stop and disables itself, so a second
// press cannot cancel something that already ended.
func (r *toolResult) cancelRun() {
	if r.cancel == nil {
		return
	}

	r.cancel()
	r.cancelButton.Disable()
	r.label.SetText("Cancelling…")
}

// progress shows how far a run has got. It is called from the worker goroutine,
// so the label is updated on the UI thread.
func (r *toolResult) progress(p desktop.Progress) {
	fyne.Do(func() {
		r.label.SetText(fmt.Sprintf("%s %d of %d — %s", r.status, p.Done+1, p.Total, p.Current))
	})
}

// detailField is one labelled value in a detail panel.
type detailField struct {
	Label string
	Value string
}

// newFieldList lays detail fields out as a read-only two-column list: a bold,
// selectable label in a column sized to the longest key, and a wrapping,
// selectable value beside it.
//
// It replaces widget.Form, whose keys are drawn trailing-aligned and bold, so a
// column of labels of different widths reads as a ragged edge next to the
// values. Here both columns start at the same left edge.
func newFieldList(fields []detailField) *fyne.Container {
	objects := make([]fyne.CanvasObject, 0, len(fields)*2)

	for _, field := range fields {
		label := widget.NewLabelWithStyle(field.Label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		label.Selectable = true

		value := widget.NewLabel(field.Value)
		value.Wrapping = fyne.TextWrapWord
		value.Selectable = true

		objects = append(objects, label, value)
	}

	return container.New(layout.NewFormLayout(), objects...)
}

// vCenterLayout keeps an object at its natural size, left-aligned and centred
// vertically in the space the row gives it.
//
// An HBox stretches every child to the row height, and a text object draws from
// its own top edge. Two different line heights — a 22 pt emoji and a 16 pt
// summary — then share a top edge and the emoji's centre lands lower, which
// reads as "the emoji is dropped".
type vCenterLayout struct{}

func (vCenterLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var size fyne.Size

	for _, object := range objects {
		if !object.Visible() {
			continue
		}

		minimum := object.MinSize()
		size.Width = max(size.Width, minimum.Width)
		size.Height = max(size.Height, minimum.Height)
	}

	return size
}

func (vCenterLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, object := range objects {
		if !object.Visible() {
			continue
		}

		minimum := object.MinSize()
		object.Resize(minimum)
		object.Move(fyne.NewPos(object.Position().X, (size.Height-minimum.Height)/2))
	}
}

// centeredVertically keeps an object at its natural size in the middle of the
// row it sits in.
func centeredVertically(object fyne.CanvasObject) *fyne.Container {
	return container.New(vCenterLayout{}, object)
}

// formRow is one label and control in a dialog or a tool card.
type formRow struct {
	label   string
	control fyne.CanvasObject
}

// newFormGrid lays a form out with the same left-aligned, selectable keys as a
// card's field list. widget.Form draws its keys as trailing-aligned rich text a
// user cannot select, which is the one thing every label in this window avoids.
func newFormGrid(rows []formRow) *fyne.Container {
	objects := make([]fyne.CanvasObject, 0, len(rows)*2)

	for _, row := range rows {
		label := widget.NewLabelWithStyle(row.label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		label.Selectable = true

		objects = append(objects, label, row.control)
	}

	return container.New(layout.NewFormLayout(), objects...)
}

// block is a message in an inset panel: an optional accent bar on the left, the
// text padded and selectable. The History raw note and a tool result both read
// as a block instead of a line of text floating in the card.
type block struct {
	accent fyne.CanvasObject
	label  *widget.Label
	panel  *fyne.Container
}

func newBlock(text string) *block {
	result := &block{label: widget.NewLabel(text)}
	result.label.Wrapping = fyne.TextWrapWord
	result.label.Selectable = true

	background := canvas.NewRectangle(theme.Color(theme.ColorNameInputBackground))
	background.CornerRadius = theme.Size(theme.SizeNameInputRadius)
	background.StrokeColor = theme.Color(theme.ColorNameSeparator)
	background.StrokeWidth = theme.Size(theme.SizeNameSeparatorThickness)

	bar := canvas.NewRectangle(theme.Color(theme.ColorNameSeparator))
	bar.SetMinSize(fyne.NewSize(3, 1))
	result.accent = bar

	body := container.NewBorder(nil, nil, bar, nil, container.NewPadded(result.label))
	result.panel = container.NewStack(background, body)

	return result
}

// setAccent colours the bar, which is what tells a reader at a glance whether a
// result succeeded, stopped or failed.
func (b *block) setAccent(name fyne.ThemeColorName) {
	if bar, ok := b.accent.(*canvas.Rectangle); ok {
		bar.FillColor = theme.Color(name)
		bar.Refresh()
	}
}

func (b *block) object() fyne.CanvasObject {
	return b.panel
}

// rowFlow places objects side by side at their natural size with a gap,
// wrapping to another row when the width runs out.
//
// An HBox cannot wrap, and a GridWrap stretches every child to one cell size, so
// neither gives spaced buttons that stay readable in a narrow window. The row
// count is remembered from the last layout, which is how MinSize reports the
// height the wrapping actually needs.
type rowFlow struct {
	gap  float32
	rows int
}

func (l *rowFlow) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var width, height float32

	for _, object := range objects {
		if object.Visible() {
			width = max(width, object.MinSize().Width)
			height = max(height, object.MinSize().Height)
		}
	}

	rows := float32(max(l.rows, 1))

	return fyne.NewSize(width, rows*height+(rows-1)*l.gap)
}

func (l *rowFlow) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	var x, y, rowHeight float32

	rows := 1

	for _, object := range objects {
		if !object.Visible() {
			continue
		}

		natural := object.MinSize()

		if x > 0 && x+natural.Width > size.Width {
			x = 0
			y += rowHeight + l.gap
			rowHeight = 0
			rows++
		}

		object.Move(fyne.NewPos(x, y))
		object.Resize(natural)
		x += natural.Width + l.gap
		rowHeight = max(rowHeight, natural.Height)
	}

	l.rows = rows
}

// spacedActions puts a card's action buttons in one row with a visible gap,
// wrapping when the card is too narrow.
func spacedActions(buttons ...fyne.CanvasObject) *fyne.Container {
	return container.New(&rowFlow{gap: theme.Padding() * 2}, buttons...)
}

// indented gives a card body the left inset that separates it from the header
// text, and a small right margin for the scrollbar.
func indented(content fyne.CanvasObject) fyne.CanvasObject {
	pad := theme.Padding()

	return container.NewBorder(
		nil,
		nil,
		gap(pad*3, 0),
		gap(pad, 0),
		content,
	)
}

// gap is an invisible spacer. An HBox's spacing is fixed, so a row that has to
// stop short of the list's scrollbar adds one of these at its end.
func gap(width, height float32) fyne.CanvasObject {
	space := canvas.NewRectangle(color.Transparent)
	space.SetMinSize(fyne.NewSize(width, height))

	return space
}

// scrollGutter is the space a viewport takes for its scrollbar plus a little
// padding. Fyne overlays the bar instead of reserving a column, so content at the
// right edge needs this much inset to stay clear of it.
func scrollGutter() float32 {
	return theme.ScrollBarSize() + theme.Padding()
}

// page pads a scrolling tab: lateral margins, and a right inset wide enough for
// the overlaid scrollbar.
func page(content fyne.CanvasObject) fyne.CanvasObject {
	pad := theme.Padding()

	return container.NewBorder(
		gap(0, pad),
		gap(0, pad),
		gap(pad*2, 0),
		gap(scrollGutter(), 0),
		content,
	)
}

// actionRow separates a card's actions from the form above it and lets them
// wrap: at a narrow width the grid stacks the buttons instead of pushing the
// last one off the card.
func actionRow(buttons ...fyne.CanvasObject) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewSeparator(),
		container.NewGridWrap(fyne.NewSize(150, 40), buttons...),
	)
}

// revealOffset is the scroll offset that brings a card into view with a small
// margin, or the current offset when it already is. It moves the minimum, so an
// expanded card stays where the user put it instead of jumping to the top.
func revealOffset(cardTop, height, viewHeight, current float32) float32 {
	margin := theme.Padding()

	if cardTop < current+margin {
		return max(0, cardTop-margin)
	}

	if bottom := cardTop + height; bottom > current+viewHeight-margin {
		return bottom - viewHeight + margin
	}

	return current
}

// newCard builds a card whose description wraps.
//
// widget.NewCard paints its subtitle with a canvas.Text, which does not wrap, so
// the subtitle's full width becomes the card's minimum and a long sentence makes
// the page wider than a narrow window (#415). The description is a wrapping
// label in the body instead.
func newCard(title, description string, content fyne.CanvasObject) fyne.CanvasObject {
	desc := widget.NewLabel(description)
	desc.Wrapping = fyne.TextWrapWord
	desc.Selectable = true

	return widget.NewCard(title, "", container.NewVBox(desc, content))
}

// newCenteredCard is newCard with its title and description centred, for the
// About page.
//
// The native card title is a left-aligned canvas.Text, so the centred card
// carries its own title label and leaves the native one empty.
func newCenteredCard(title, description string, content fyne.CanvasObject) fyne.CanvasObject {
	head := widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	desc := widget.NewLabel(description)
	desc.Wrapping = fyne.TextWrapWord
	desc.Alignment = fyne.TextAlignCenter
	desc.Selectable = true

	return widget.NewCard("", "", container.NewVBox(head, desc, content))
}

// appendPath adds a picked path to a comma-separated list, so repeated Browse
// clicks build a multi-path value without a multi-select dialog.
func appendPath(target *widget.Entry, path string) {
	current := strings.TrimSpace(target.Text)
	if current == "" {
		target.SetText(path)
		return
	}

	target.SetText(current + ", " + path)
}

// uriPath is the file-system path of a Fyne URI. On Windows Fyne reports
// /C:/Users/... — the leading slash would make the CLI see a path that does not
// exist, so it is dropped.
func uriPath(uri fyne.URI) string {
	path := filepath.FromSlash(uri.Path())

	if len(path) > 2 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}

	return path
}

// browseFiles appends one picked file to a field. Fyne's file dialog is
// single-select, so picking several means pressing Browse again; the field keeps
// the earlier picks.
func (u *ui) browseFiles(target *widget.Entry) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}

		_ = reader.Close()
		appendPath(target, uriPath(reader.URI()))
	}, u.window)
}

// browseFolder fills a field with one picked folder.
func (u *ui) browseFolder(target *widget.Entry) {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil || uri == nil {
			return
		}

		target.SetText(uriPath(uri))
	}, u.window)
}
