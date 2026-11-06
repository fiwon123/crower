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

// Spaced actions lay their buttons out in one row with a gap, and wrap instead
// of running past the card.
func TestSpacedActionsWrapAndKeepAGap(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	buttons := []fyne.CanvasObject{
		widget.NewButton("Run", nil),
		widget.NewButton("Edit", nil),
		widget.NewButton("Delete", nil),
	}

	flow := &rowFlow{gap: theme.Padding() * 2}

	// Widths are measured, not assumed: a platform's font decides how wide a
	// button is, and the wrap point has to be stated in those terms.
	var wide float32
	for _, button := range buttons {
		wide += button.MinSize().Width + theme.Padding()*2
	}

	flow.Layout(buttons, fyne.NewSize(wide, 40))

	for index := 1; index < len(buttons); index++ {
		previous := buttons[index-1]
		current := buttons[index]
		gap := current.Position().X - (previous.Position().X + previous.Size().Width)

		if gap != theme.Padding()*2 {
			t.Errorf("gap between buttons = %v, want %v", gap, theme.Padding()*2)
		}
	}

	narrow := buttons[0].MinSize().Width + theme.Padding()*2 + buttons[1].MinSize().Width - 1
	flow.Layout(buttons, fyne.NewSize(narrow, 40))

	if buttons[1].Position().Y <= buttons[0].Position().Y {
		t.Error("a narrow row did not wrap its buttons")
	}

	if size := flow.MinSize(buttons); size.Height < buttons[0].Size().Height*2 {
		t.Errorf("wrapped minimum height = %v, want room for two rows", size.Height)
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
