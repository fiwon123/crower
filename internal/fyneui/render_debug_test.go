//go:build fyne

package fyneui

import (
	"fmt"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// findLabels walks a rendered tree and collects every label with its geometry.
func findLabels(object fyne.CanvasObject, out *[]string) {
	switch typed := object.(type) {
	case *widget.Label:
		*out = append(*out, fmt.Sprintf("label %q pos=%v size=%v visible=%v", typed.Text, typed.Position(), typed.Size(), typed.Visible()))
	case *fyne.Container:
		for _, child := range typed.Objects {
			findLabels(child, out)
		}
	case fyne.Widget:
		if r := test.WidgetRenderer(typed); r != nil {
			for _, child := range r.Objects() {
				findLabels(child, out)
			}
		}
	}
}

func TestDiagnosticRowText(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})

	canvas := test.NewCanvasWithPainter(software.NewPainter())
	canvas.SetContent(front.content())
	canvas.Resize(fyne.NewSize(960, 640))

	labels := make([]string, 0)
	findLabels(canvas.Content(), &labels)

	img := canvas.Capture()
	bounds := img.Bounds()

	var ink strings.Builder

	for y := bounds.Min.Y; y < 320 && y < bounds.Max.Y; y += 8 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 8 {
			colours := map[string]bool{}

			for dy := 0; dy < 8 && y+dy < bounds.Max.Y; dy++ {
				for dx := 0; dx < 8 && x+dx < bounds.Max.X; dx++ {
					r, g, b, a := img.At(x+dx, y+dy).RGBA()
					colours[fmt.Sprintf("%d,%d,%d,%d", r, g, b, a)] = true
				}
			}

			if len(colours) > 1 {
				ink.WriteByte('#')
			} else {
				ink.WriteByte('.')
			}
		}

		ink.WriteByte('\n')
	}

	t.Errorf("canvas=%v list=%v\nlabels:\n%s\nink:\n%s", bounds, front.list.Size(), strings.Join(labels, "\n"), ink.String())
}
