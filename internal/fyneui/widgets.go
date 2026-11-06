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
	result := &toolResult{label: widget.NewLabel("")}
	result.label.Wrapping = fyne.TextWrapWord

	result.cancelButton = widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), result.cancelRun)
	result.cancelButton.Importance = widget.LowImportance
	result.cancelButton.Hide()

	result.body = container.NewVBox(result.label, result.cancelButton)
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

// report shows an outcome and opens the area: a result square appears because an
// action ran, not before.
func (r *toolResult) report(text string) {
	r.label.SetText(text)
	r.area.Show()
	r.setVisible(true)
}

// start opens the area for a run and arms the cancel action. A nil cancel means
// the run cannot be stopped, so no Cancel button is offered.
func (r *toolResult) start(status string, cancel context.CancelFunc) {
	r.status = status
	r.cancel = cancel

	r.label.SetText(status)
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

// setFieldGrid replaces a two-column grid with label/value rows. Keys are bold,
// values wrap, and an empty list leaves the grid empty — both detail panels use
// it, so a command and a history entry read the same way.
func setFieldGrid(grid *fyne.Container, fields []detailField) {
	if grid == nil {
		return
	}

	grid.RemoveAll()

	for _, field := range fields {
		label := widget.NewLabelWithStyle(field.Label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		label.Wrapping = fyne.TextWrapWord

		value := widget.NewLabel(field.Value)
		value.Wrapping = fyne.TextWrapWord

		grid.Add(label)
		grid.Add(value)
	}

	grid.Refresh()
}

// gap is an invisible spacer. An HBox's spacing is fixed, so a row that has to
// stop short of the list's scrollbar adds one of these at its end.
func gap(width, height float32) fyne.CanvasObject {
	space := canvas.NewRectangle(color.Transparent)
	space.SetMinSize(fyne.NewSize(width, height))

	return space
}

// page pads a scrolling tab. Fyne's scrollbar overlays the viewport instead of
// reserving a column, so the right inset is wider than the edge padding: without
// it a control at the edge sits under the scrollbar.
func page(content fyne.CanvasObject) fyne.CanvasObject {
	pad := theme.Padding()

	return container.NewBorder(
		gap(0, pad),
		gap(0, pad),
		gap(pad*2, 0),
		gap(pad*3, 0),
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
