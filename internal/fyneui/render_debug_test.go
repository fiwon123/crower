//go:build fyne

package fyneui

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/test"
)

// distinctColors counts the colours in a region of the captured canvas. Text
// antialiasing produces several; a blank area produces one or two.
func distinctColors(img interface{ At(x, y int) color.Color }, x0, y0, x1, y1 int) int {
	seen := map[color.RGBA64]bool{}

	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			seen[color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)}] = true
		}
	}

	return len(seen)
}

// Scratch diagnostic: is there any ink where the first command row is drawn?
func TestCommandRowPaintsInk(t *testing.T) {
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	window.SetContent(front.content())
	window.Resize(fyneSize())

	img := test.Canvas().Capture()

	rows := distinctColors(img, 20, 100, 500, 150)
	t.Logf("command row band colours: %d", rows)

	if rows < 4 {
		t.Errorf("command row region has %d colours; the rows look blank", rows)
	}

	front.tabs.SelectIndex(1)

	img = test.Canvas().Capture()

	history := distinctColors(img, 20, 100, 500, 150)
	t.Logf("history band colours: %d", history)

	if history < 4 {
		t.Errorf("history region has %d colours; the table looks blank", history)
	}
}
