//go:build fyne

package fyneui

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// growingRows counts explicit lines and wrapped rows, between the bounds the
// field was built with.
func TestGrowingRows(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		minRows int
		maxRows int
		want    int
	}{
		{name: "empty", text: "", minRows: 4, maxRows: 12, want: 4},
		{name: "one short line", text: "echo hi", minRows: 4, maxRows: 12, want: 4},
		{name: "six lines", text: "a\nb\nc\nd\ne\nf", minRows: 4, maxRows: 12, want: 6},
		{name: "capped", text: strings.Repeat("line\n", 20), minRows: 4, maxRows: 12, want: 12},
		{name: "one wrapped line", text: strings.Repeat("x", 150), minRows: 4, maxRows: 12, want: 4},
		{name: "one long wrapped line", text: strings.Repeat("x", 300), minRows: 4, maxRows: 12, want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := growingRows(tt.text, tt.minRows, tt.maxRows); got != tt.want {
				t.Errorf("growingRows(%q) = %d, want %d", tt.text, got, tt.want)
			}
		})
	}
}

// Setting text on a growing field raises its minimum height, which is what the
// dialog lays out; the cap keeps it from taking the window.
func TestGrowingEntryGrowsWithItsText(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	entry := newGrowingEntry(4, 10)
	entry.Resize(fyne.NewSize(600, 120))

	short := entry.MinSize().Height

	entry.SetText("1\n2\n3\n4\n5\n6\n7\n8")

	tall := entry.MinSize().Height

	if tall <= short {
		t.Errorf("MinSize height = %v, want it to grow past %v", tall, short)
	}

	entry.SetText(strings.Repeat("line\n", 40))

	capped := entry.MinSize().Height

	if got := growingRows(entry.Text, 4, 10); got != 10 {
		t.Fatalf("growingRows = %d, want the cap", got)
	}

	if capped < tall {
		t.Errorf("capped MinSize height = %v, want at least %v", capped, tall)
	}
}

// A detail list starts both columns at the left edge and keeps every label and
// value selectable.
func TestFieldListAlignsLabelsAndKeepsTextSelectable(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	list := newFieldList([]detailField{
		{Label: "Exec", Value: "echo hi"},
		{Label: "Shell", Value: "sh"},
	})

	if len(list.Objects) != 4 {
		t.Fatalf("list = %d objects, want two label/value pairs", len(list.Objects))
	}

	for index, want := range []string{"Exec", "echo hi", "Shell", "sh"} {
		label, ok := list.Objects[index].(*widget.Label)
		if !ok {
			t.Fatalf("object %d = %T, want a label", index, list.Objects[index])
		}

		if label.Text != want {
			t.Errorf("object %d = %q, want %q", index, label.Text, want)
		}

		if !label.Selectable {
			t.Errorf("%q is not selectable", label.Text)
		}
	}

	for _, index := range []int{0, 2} {
		label := list.Objects[index].(*widget.Label)
		if label.Alignment != fyne.TextAlignLeading {
			t.Errorf("label %q is not left-aligned", label.Text)
		}

		if !label.TextStyle.Bold {
			t.Errorf("label %q is not emphasised", label.Text)
		}
	}
}

// Spaced actions lay their buttons out in one row with the shared small gap, and
// wrap instead of running past the card.
func TestSpacedActionsWrapAndKeepAGap(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	buttons := []fyne.CanvasObject{
		widget.NewButton("Run", nil),
		widget.NewButton("Edit", nil),
		widget.NewButton("Delete", nil),
	}

	// The product's own row, so the test guards the gap it ships.
	row := spacedActions(buttons...)

	// Widths are measured, not assumed: a platform's font decides how wide a
	// button is, and the wrap point has to be stated in those terms.
	var wide float32
	for _, button := range buttons {
		wide += button.MinSize().Width + actionGap()
	}

	row.Resize(fyne.NewSize(wide, 40))

	for index := 1; index < len(buttons); index++ {
		previous := buttons[index-1]
		current := buttons[index]
		gap := current.Position().X - (previous.Position().X + previous.Size().Width)

		if gap != actionGap() {
			t.Errorf("gap between buttons = %v, want %v", gap, actionGap())
		}
	}

	narrow := buttons[0].MinSize().Width + actionGap() + buttons[1].MinSize().Width - 1
	row.Resize(fyne.NewSize(narrow, 40))

	if buttons[1].Position().Y <= buttons[0].Position().Y {
		t.Error("a narrow row did not wrap its buttons")
	}

	if size := row.MinSize(); size.Height < buttons[0].Size().Height*2 {
		t.Errorf("wrapped minimum height = %v, want room for two rows", size.Height)
	}
}

// Every action button shares one height, whatever its icon or label, and an
// icon-only button stays square.
func TestNewButtonSharesOneSize(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	withIcon := newButton("Run", theme.MediaPlayIcon(), widget.SuccessImportance, nil)
	textOnly := newButton("Cancel", nil, widget.LowImportance, nil)
	iconOnly := newButton("", theme.DeleteIcon(), widget.DangerImportance, nil)

	for name, button := range map[string]*appButton{
		"with icon": withIcon,
		"text only": textOnly,
		"icon only": iconOnly,
	} {
		if got := button.MinSize().Height; got != buttonHeight {
			t.Errorf("%s height = %.1f, want %.1f", name, got, buttonHeight)
		}
	}

	for name, button := range map[string]*appButton{"with icon": withIcon, "text only": textOnly} {
		if got, want := button.MinSize().Width, button.Button.MinSize().Width; got != want {
			t.Errorf("%s width = %.1f, want the natural width %.1f", name, got, want)
		}
	}

	if got := iconOnly.MinSize().Width; got != buttonHeight {
		t.Errorf("icon-only width = %.1f, want %.1f", got, buttonHeight)
	}

	if withIcon.Importance != widget.SuccessImportance {
		t.Errorf("importance = %v, want success", withIcon.Importance)
	}
}

// The embedded button is what the renderer paints, so a later SetText or
// SetIcon repaints the button on screen rather than a second unseen copy.
func TestAppButtonPaintsItsEmbeddedButton(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	button := newButton("Install", theme.DownloadIcon(), widget.HighImportance, nil)

	renderer := button.CreateRenderer()
	if len(renderer.Objects()) != 1 {
		t.Fatalf("renderer has %d objects, want the embedded button", len(renderer.Objects()))
	}

	if renderer.Objects()[0] != fyne.CanvasObject(button.Button) {
		t.Fatalf("renderer paints %T, want the embedded button", renderer.Objects()[0])
	}

	button.SetText("Install v1.1.0")

	if button.Button.Text != "Install v1.1.0" {
		t.Errorf("embedded text = %q, want the new label", button.Button.Text)
	}
}

// A container positions the wrapper, not the button inside it. The promoted
// geometry used to move the embedded button, which the renderer then painted at
// that offset a second time — so every button after the first drifted right
// (#507).
func TestAppButtonPaintsWhereTheContainerPutIt(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	first := newButton("Cancel", nil, widget.LowImportance, nil)
	second := newButton("Create", nil, widget.HighImportance, nil)

	window := test.NewWindow(spacedActions(first, second))
	t.Cleanup(window.Close)
	window.Resize(fyne.NewSize(400, buttonHeight))

	driver := fyne.CurrentApp().Driver()

	for _, button := range []*appButton{first, second} {
		outer := driver.AbsolutePositionForObject(button)
		inner := driver.AbsolutePositionForObject(button.Button)

		if inner != outer {
			t.Errorf("%q is painted at x=%.1f but positioned at x=%.1f: the wrapped button moved",
				button.Text, inner.X, outer.X)
		}
	}

	firstPos := driver.AbsolutePositionForObject(first.Button)
	secondPos := driver.AbsolutePositionForObject(second.Button)

	if gap := secondPos.X - (firstPos.X + first.Button.Size().Width); gap != actionGap() {
		t.Errorf("painted gap = %.1f, want %v", gap, actionGap())
	}
}

// buttonsIn collects every action button the canvas can reach, walking the tree
// the way Fyne paints it.
func buttonsIn(object fyne.CanvasObject) []*appButton {
	switch typed := object.(type) {
	case *appButton:
		return []*appButton{typed}
	case *fyne.Container:
		found := make([]*appButton, 0, len(typed.Objects))

		for _, child := range typed.Objects {
			found = append(found, buttonsIn(child)...)
		}

		return found
	case fyne.Widget:
		renderer := test.WidgetRenderer(typed)
		if renderer == nil {
			return nil
		}

		found := make([]*appButton, 0, len(renderer.Objects()))

		for _, child := range renderer.Objects() {
			found = append(found, buttonsIn(child)...)
		}

		return found
	}

	return nil
}

// The built window really uses the one constructor everywhere: no button in any
// tab is taller or shorter than another (#492).
func TestEveryWindowButtonSharesOneHeight(t *testing.T) {
	front, window := newTestUI(t, testModel(), stubEditor{}, stubReverter{})

	window.SetContent(front.content())
	window.Resize(fyne.NewSize(960, 640))

	// AppTabs attaches only the selected tab, so walk each tab's content
	// directly. A command card builds its actions on first open.
	for _, card := range front.commandCards {
		if card.key == "greet" {
			front.toggleCommandCard(card)
			break
		}
	}

	buttons := make([]*appButton, 0, 32)

	for _, item := range front.tabs.Items {
		buttons = append(buttons, buttonsIn(item.Content)...)
	}

	if len(buttons) < 20 {
		t.Fatalf("found %d buttons, want the tabs' actions", len(buttons))
	}

	for _, button := range buttons {
		if got := button.MinSize().Height; got != buttonHeight {
			t.Errorf("button %q height = %.1f, want %.1f", button.Text, got, buttonHeight)
		}
	}
}

// A dialog footer keeps its buttons as one centred row: the native Fyne dialog
// put them at the right edge, which is the inconsistency #500 removes.
func TestDialogFooterCentersItsButtons(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	confirm := newButton("Yes", nil, widget.HighImportance, nil)
	dismiss := newButton("No", nil, widget.LowImportance, nil)

	window := test.NewWindow(dialogFooter(dismiss, confirm))
	window.Resize(fyne.NewSize(600, 120))

	buttons := buttonsIn(window.Canvas().Content())
	if len(buttons) != 2 {
		t.Fatalf("found %d buttons, want the footer's two", len(buttons))
	}

	driver := fyne.CurrentApp().Driver()

	var left, right float32

	for index, button := range buttons {
		x := driver.AbsolutePositionForObject(button).X
		if index == 0 || x < left {
			left = x
		}

		if edge := x + button.Size().Width; index == 0 || edge > right {
			right = edge
		}
	}

	if center, want := left+(right-left)/2, float32(300); center != want {
		t.Errorf("button row centre = %.1f, want %.1f (left %.1f, right %.1f)", center, want, left, right)
	}
}

// The About action stack uses a gap tighter than the theme's padding.
func TestVStackUsesTheTighterGap(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	first := newButton("One", nil, widget.MediumImportance, nil)
	second := newButton("Two", nil, widget.MediumImportance, nil)

	stack := &vStack{gap: theme.Padding() / 2, center: true}
	objects := []fyne.CanvasObject{first, second}

	if got, want := stack.MinSize(objects).Height, first.MinSize().Height*2+stack.gap; got != want {
		t.Errorf("height = %.1f, want %.1f", got, want)
	}

	stack.Layout(objects, fyne.NewSize(400, stack.MinSize(objects).Height))

	if got, want := second.Position().Y, first.MinSize().Height+stack.gap; got != want {
		t.Errorf("second row y = %.1f, want %.1f", got, want)
	}

	if got, want := first.Position().X, (400-first.MinSize().Width)/2; got != want {
		t.Errorf("first row x = %.1f, want %.1f", got, want)
	}
}

// A result block colours its accent bar by the outcome.
func TestBlockAccentFollowsTheOutcome(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	result := newToolResult()
	result.report("done")

	bar, ok := result.panel.accent.(*canvas.Rectangle)
	if !ok {
		t.Fatalf("accent = %T, want a rectangle", result.panel.accent)
	}

	if bar.FillColor != theme.Color(theme.ColorNameSuccess) {
		t.Errorf("success accent = %v, want the success colour", bar.FillColor)
	}

	result.reportWith("error: permission denied", theme.ColorNameError)

	if bar.FillColor != theme.Color(theme.ColorNameError) {
		t.Errorf("failure accent = %v, want the error colour", bar.FillColor)
	}
}

// Output and error text can be selected and copied.
func TestOutputTextIsSelectable(t *testing.T) {
	label := newSelectableText("command failed")
	if !label.Selectable {
		t.Error("the output label is not selectable")
	}
}

// A tool result's outcome can be selected and copied too.
func TestToolResultIsSelectable(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	if result := newToolResult(); !result.label.Selectable {
		t.Error("the tool result label is not selectable")
	}
}

// A card's description is selectable like the values.
func TestCardDescriptionsAreSelectable(t *testing.T) {
	card := newCard("Title", "description", widget.NewLabel("content"))

	label := findLabel(card, "description")
	if label == nil || !label.Selectable {
		t.Error("the card description is missing or not selectable")
	}
}

// Revealing a card scrolls the minimum: visible stays put, above and below move
// just enough.
func TestRevealOffset(t *testing.T) {
	tests := []struct {
		name    string
		top     float32
		height  float32
		view    float32
		current float32
		want    float32
	}{
		{"already visible", 100, 50, 400, 0, 0},
		{"above the viewport", 10, 50, 400, 100, 6},
		{"below the viewport", 500, 80, 400, 0, 184},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := revealOffset(tt.top, tt.height, tt.view, tt.current); got != tt.want {
				t.Errorf("revealOffset = %.1f, want %.1f", got, tt.want)
			}
		})
	}
}
