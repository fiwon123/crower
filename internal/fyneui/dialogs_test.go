//go:build fyne

package fyneui

import (
	"sync/atomic"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
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

// A destructive action stays disabled until its name is typed exactly.
func TestNamedConfirmationNeedsTheExactName(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	window := test.NewWindow(nil)
	t.Cleanup(window.Close)

	var confirmed atomic.Bool

	pop := confirmNamed("Delete command", `Delete "greet"?`, "greet", "Delete", widget.DangerImportance, window, func(ok bool) {
		if ok {
			confirmed.Store(true)
		}
	})

	action := findButton(pop, "Delete")
	if action == nil {
		t.Fatal("the dialog has no Delete button")
	}

	if !action.Disabled() {
		t.Error("the Delete button starts enabled")
	}

	entry := entryIn(pop)
	if entry == nil {
		t.Fatal("the dialog has no name field")
	}

	entry.SetText("gre")

	if !action.Disabled() {
		t.Error("a partial name enabled the Delete button")
	}

	entry.SetText("greet")

	if action.Disabled() {
		t.Fatal("the exact name did not enable the Delete button")
	}

	action.OnTapped()

	if !confirmed.Load() {
		t.Error("the enabled action did not run")
	}
}

// entryIn returns the first editable entry inside an object.
func entryIn(object fyne.CanvasObject) *widget.Entry {
	switch typed := object.(type) {
	case *widget.Entry:
		return typed
	case *fyne.Container:
		for _, child := range typed.Objects {
			if entry := entryIn(child); entry != nil {
				return entry
			}
		}
	case fyne.Widget:
		renderer := test.WidgetRenderer(typed)
		if renderer == nil {
			return nil
		}

		for _, child := range renderer.Objects() {
			if entry := entryIn(child); entry != nil {
				return entry
			}
		}
	}

	return nil
}
