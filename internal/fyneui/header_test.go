//go:build fyne

package fyneui

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
)

// The emoji and the summary are one text run: a dropdown aligns them that way,
// and two objects of different sizes never share a baseline (#480).
func TestHeaderDrawsOneLine(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	header := newCardHeader("🚀", "deploy")

	if header.text.Text != "🚀  deploy" {
		t.Errorf("header line = %q, want the emoji and the summary together", header.text.Text)
	}

	if header.text.TextSize < 18 {
		t.Errorf("header text is %.1f pt, want it readable", header.text.TextSize)
	}

	bare := newCardHeader("", "deploy")

	if bare.text.Text != "deploy" {
		t.Errorf("header without an icon = %q, want the bare summary", bare.text.Text)
	}
}

// The centring layout keeps an object at its natural size and moves it to the
// middle of the row, leaving its horizontal position alone.
func TestVCenterLayoutCentersVertically(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	object := canvas.NewText("x", color.White)

	layout := vCenterLayout{}
	layout.Layout([]fyne.CanvasObject{object}, fyne.NewSize(100, 40))

	want := (40 - object.MinSize().Height) / 2
	if got := object.Position().Y; got != want {
		t.Errorf("y = %.1f, want the object centred in a 40 px row (%.1f)", got, want)
	}

	if got := object.Position().X; got != 0 {
		t.Errorf("x = %.1f, want the object left-aligned", got)
	}
}
