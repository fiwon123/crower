//go:build fyne

package fyneui

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

// The emoji and the summary are centred on the same line: an HBox alone left the
// taller glyph's centre below the text's (#477).
func TestHeaderCentersTheEmojiOnTheSummary(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	header := newCardHeader("🚀", "deploy")

	window := test.NewWindow(container.NewVBox(header))
	t.Cleanup(window.Close)

	window.Resize(fyne.NewSize(420, 120))

	driver := fyne.CurrentApp().Driver()

	emojiAt := driver.AbsolutePositionForObject(header.emoji)
	titleAt := driver.AbsolutePositionForObject(header.title)

	emojiCenter := emojiAt.Y + header.emoji.Size().Height/2
	titleCenter := titleAt.Y + header.title.Size().Height/2

	if difference := emojiCenter - titleCenter; difference > 1 || difference < -1 {
		t.Errorf("emoji centre = %.1f, summary centre = %.1f, want them on one line", emojiCenter, titleCenter)
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
