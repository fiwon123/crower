//go:build fyne

package fyneui

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

// Output and error dialogs are non-modal pop-ups: they can be dismissed, and a
// click outside closes them.
func TestOutputDialogIsDismissible(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	window := test.NewWindow(nil)
	t.Cleanup(window.Close)

	pop := showText("Output", "command failed", window)

	if len(window.Canvas().Overlays().List()) == 0 {
		t.Fatal("showText did not open an overlay")
	}

	pop.Hide()

	if len(window.Canvas().Overlays().List()) != 0 {
		t.Error("hiding the output left the overlay behind")
	}
}
