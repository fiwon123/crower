//go:build fyne

package fyneui

import (
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// newScrollFreeEntry builds a single-line field that does not swallow the mouse
// wheel.
//
// Fyne dispatches a wheel event to the topmost Scrollable under the cursor and
// never bubbles it to a parent (internal/driver/glfw window.processMouseScrolled).
// Every widget.Entry builds an internal scroll unless its wrapping is
// TextWrapOff and its Scroll is ScrollNone (widget/entry.go), and the
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
	entry.Scroll = widget.ScrollNone

	return entry
}

// toolResult is one card's own result area. It stays hidden until the card runs
// something, then opens with the outcome; the toggle hides it again, which is
// what gives the forms back their space.
type toolResult struct {
	label  *widget.Label
	toggle *widget.Button
	body   *fyne.Container
	area   *fyne.Container
	shown  bool
}

func newToolResult() *toolResult {
	result := &toolResult{label: widget.NewLabel("")}
	result.label.Wrapping = fyne.TextWrapWord

	result.body = container.NewVBox(result.label)
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
