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

// Scratch diagnostic: where is any ink drawn, and is there any at the row band?
func TestCommandRowPaintsInk(t *testing.T) {
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	window.SetContent(front.content())
	window.Resize(fyneSize())

	img := test.Canvas().Capture()
	bounds := img.Bounds()
	t.Logf("capture bounds: %v", bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 40 {
		colours := distinctColors(img, 0, y, bounds.Max.X, min(y+40, bounds.Max.Y))
		t.Logf("band y=%d..%d colours=%d", y, min(y+40, bounds.Max.Y), colours)
	}

	t.Logf("list size: %v, visible: %v", front.list.Size(), front.list.Visible())
	t.Logf("first row text: %q", front.model.FilteredKeys())
}
