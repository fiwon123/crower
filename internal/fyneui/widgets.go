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
	toggle       *appButton
	cancelButton *appButton
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

	result.cancelButton = newButton("Cancel", theme.CancelIcon(), widget.LowImportance, result.cancelRun)
	result.cancelButton.Hide()

	result.body = container.NewVBox(result.panel.object(), result.cancelButton)
	result.body.Hide()

	result.toggle = newButton("Hide result", nil, widget.LowImportance, func() {
		result.setVisible(!result.shown)
	})

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

// gutterLayout places one object after a left gutter and centres it vertically
// in the row.
//
// A Border already pads its right slot by theme.Padding(), and a box layout pads
// between its children, so either would add more than the gutter asked for; a
// single-child box adds nothing. This layout states both the gap and the centring
// exactly (#510).
type gutterLayout struct {
	left float32
}

func (l gutterLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var size fyne.Size

	for _, object := range objects {
		if !object.Visible() {
			continue
		}

		minimum := object.MinSize()
		size.Width = max(size.Width, minimum.Width)
		size.Height = max(size.Height, minimum.Height)
	}

	size.Width += l.left

	return size
}

func (l gutterLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, object := range objects {
		if !object.Visible() {
			continue
		}

		minimum := object.MinSize()
		object.Resize(minimum)
		object.Move(fyne.NewPos(l.left, (size.Height-minimum.Height)/2))
	}
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

// The window's one button height. Fyne sizes a button from its content — a text
// label is shorter than an inline icon — so an icon button and a text-only
// button side by side differ in height by 8 px, and a GridWrap cell overrides
// both. Every action button goes through newButton so it shares this height.
//
// Width is the opposite rule: a button takes the width its label needs. A
// minimum width stretched short labels ("Run", "Edit") and pushed the row apart,
// and it inflated the About and dialog buttons, so there is none (#495).
const buttonHeight float32 = 40

// appButton is a widget.Button with the window's single button height. The
// embedded button keeps its renderer and handlers; only the geometry belongs to
// this wrapper.
//
// The wrapper owns its position and size deliberately. A container lays out the
// object it was given, so Move and Resize arrive here; promoted to the embedded
// button they would move the button that CreateRenderer paints, and the renderer
// would then paint it at that offset a second time. The first child of a row sat
// still at x=0 while every later one drifted right — the dialogs' second button,
// the About column and every action row (#507).
type appButton struct {
	*widget.Button

	iconOnly bool

	pos  fyne.Position
	size fyne.Size
}

// MinSize reports the canonical height and the width the embedded button asks
// for. An icon-only button is square, so a row of them lines up with a row of
// labels.
func (b *appButton) MinSize() fyne.Size {
	if b.iconOnly {
		return fyne.NewSize(buttonHeight, buttonHeight)
	}

	return fyne.NewSize(b.Button.MinSize().Width, buttonHeight)
}

// CreateRenderer paints the embedded button through this widget. Without it the
// canvas would render a second copy, and a later SetText or SetIcon on the
// embedded button would repaint only the copy nobody can see.
func (b *appButton) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(b.Button)
}

// Move places the wrapper. The embedded button stays at the origin, so the
// renderer paints it exactly once, at its container's position.
func (b *appButton) Move(pos fyne.Position) {
	b.pos = pos
}

// Position is where the container put this wrapper.
func (b *appButton) Position() fyne.Position {
	return b.pos
}

// Resize sizes the wrapper and the button it paints.
func (b *appButton) Resize(size fyne.Size) {
	b.size = size
	b.Button.Resize(size)
}

// Size is the size the container gave this wrapper.
func (b *appButton) Size() fyne.Size {
	return b.size
}

// newButton builds every action button in the window: one height, one minimum
// width, and the importance passed in. A nil icon means a text-only button; an
// empty label with an icon means an icon-only button.
func newButton(text string, icon fyne.Resource, importance widget.Importance, tapped func()) *appButton {
	var button *widget.Button
	if icon == nil {
		button = widget.NewButton(text, tapped)
	} else {
		button = widget.NewButtonWithIcon(text, icon, tapped)
	}

	button.Importance = importance

	return &appButton{Button: button, iconOnly: strings.TrimSpace(text) == "" && icon != nil}
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

// vStack lays objects out vertically with a gap tighter than the theme's
// padding. The About actions use it so a short list of update steps reads as
// one group instead of four unrelated rows (#500).
type vStack struct {
	gap    float32
	center bool
}

func (l *vStack) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var width, height float32

	visible := 0

	for _, object := range objects {
		if !object.Visible() {
			continue
		}

		natural := object.MinSize()
		width = max(width, natural.Width)
		height += natural.Height
		visible++
	}

	if visible > 1 {
		height += l.gap * float32(visible-1)
	}

	return fyne.NewSize(width, height)
}

func (l *vStack) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	var y float32

	for _, object := range objects {
		if !object.Visible() {
			continue
		}

		natural := object.MinSize()

		x := float32(0)
		if l.center {
			x = max(0, (size.Width-natural.Width)/2)
		}

		object.Move(fyne.NewPos(x, y))
		object.Resize(natural)
		y += natural.Height + l.gap
	}
}

// actionGap is the space between two of a row's action buttons: the base
// spacing unit, the same gap the theme puts between any two related controls.
// A rowFlow needs it as a number, and the Browse gutter uses it so a field and
// its button are separated by exactly the same amount (#516).
func actionGap() float32 {
	return theme.Padding()
}

// spacedActions puts a card's action buttons in one row with that gap, wrapping
// when the card is too narrow.
func spacedActions(buttons ...fyne.CanvasObject) *fyne.Container {
	return container.New(&rowFlow{gap: actionGap()}, buttons...)
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

// slot is an invisible border slot that leaves exactly width and height of
// clear space beside the content. A Border adds its own base-unit padding
// around every slot, so the spacer carries the size minus that padding;
// without the subtraction a page margin came out a base unit larger than the
// value it was given (#516).
func slot(width, height float32) fyne.CanvasObject {
	pad := theme.Padding()

	return gap(max(0, width-pad), max(0, height-pad))
}

// page pads a scrolling tab: a section gap from the window's edges and a right
// inset wide enough for the overlaid scrollbar. The top toolbar uses the same
// slots, so its controls line up with the cards below (#516).
func page(content fyne.CanvasObject) fyne.CanvasObject {
	return container.NewBorder(
		slot(0, appSection),
		slot(0, appSection),
		slot(appSection, 0),
		slot(scrollGutter(), 0),
		content,
	)
}

// actionRow separates a card's actions from the form above it and lets them
// wrap: at a narrow width the grid stacks the buttons instead of pushing the
// last one off the card. It shares spacedActions' layout, so a tool card's
// buttons look like a command card's.
func actionRow(buttons ...fyne.CanvasObject) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewSeparator(),
		spacedActions(buttons...),
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

// newCard builds one of the window's standard cards: a bold title, a wrapping
// selectable description and the content, on the same inset panel the command
// and history cards use.
//
// widget.NewCard was replaced rather than wrapped: it indents its title further
// than its content and draws the title with a left-aligned canvas text, so a
// Tool card's text started on a different column from every other card in the
// window (#415, #516).
func newCard(title, description string, content fyne.CanvasObject) fyne.CanvasObject {
	return newAppCard(title, description, content, false)
}

// newCenteredCard is newCard with its title and description centred, for the
// About page.
func newCenteredCard(title, description string, content fyne.CanvasObject) fyne.CanvasObject {
	return newAppCard(title, description, content, true)
}

// newAppCard is the one card treatment. The title is a canvas.Text so it can
// share the size and weight of an expandable card's title; a widget.Label only
// takes the body size. The description stays a label so it wraps and can be
// selected.
func newAppCard(title, description string, content fyne.CanvasObject, centered bool) fyne.CanvasObject {
	head := canvas.NewText(title, theme.Color(theme.ColorNameForeground))
	head.TextSize = cardTextSize
	head.TextStyle = fyne.TextStyle{Bold: true}

	desc := widget.NewLabel(description)
	desc.Wrapping = fyne.TextWrapWord
	desc.Selectable = true

	if centered {
		head.Alignment = fyne.TextAlignCenter
		desc.Alignment = fyne.TextAlignCenter
	}

	// Title and description are one heading block; the body is a section below
	// it, the same split an expandable card makes between its header and body.
	heading := container.New(layout.NewCustomPaddedVBoxLayout(appPadding), head, desc)
	inner := container.New(layout.NewCustomPaddedVBoxLayout(appSection), heading, content)

	background := canvas.NewRectangle(theme.Color(theme.ColorNameInputBackground))
	background.CornerRadius = theme.Size(theme.SizeNameInputRadius)
	background.StrokeColor = theme.Color(theme.ColorNameSeparator)
	background.StrokeWidth = theme.Size(theme.SizeNameSeparatorThickness)

	return container.NewStack(background, container.NewPadded(inner))
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
