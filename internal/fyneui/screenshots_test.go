//go:build fyne

package fyneui

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// dumpButtons logs every button and select reachable from object, with its
// absolute position and size, so a layout can be judged from the CI log.
func dumpButtons(t *testing.T, label string, object fyne.CanvasObject) {
	t.Helper()

	t.Logf("=== %s ===", label)

	var walk func(fyne.CanvasObject)
	walk = func(o fyne.CanvasObject) {
		switch typed := o.(type) {
		case *appButton:
			pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(typed)
			t.Logf("  btn  %-22q abs=(%.0f,%.0f) size=(%.0f,%.0f)", typed.Text, pos.X, pos.Y, typed.Size().Width, typed.Size().Height)
		case *widget.Select:
			pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(typed)
			t.Logf("  sel  %-22q abs=(%.0f,%.0f) size=(%.0f,%.0f)", typed.Selected, pos.X, pos.Y, typed.Size().Width, typed.Size().Height)
		case *fyne.Container:
			for _, child := range typed.Objects {
				walk(child)
			}
		case fyne.Widget:
			if renderer := test.WidgetRenderer(typed); renderer != nil {
				for _, child := range renderer.Objects() {
					walk(child)
				}
			}
		}
	}

	walk(object)
}

// captureShot renders content to a PNG under build/shots with the software
// painter, so a layout can be inspected without opening a window.
func captureShot(t *testing.T, content fyne.CanvasObject, name string, size fyne.Size) {
	t.Helper()

	canvas := software.NewCanvas()
	canvas.SetContent(content)
	canvas.Resize(size)

	img := canvas.Capture()

	path := filepath.Join("..", "..", "build", "shots", name+".png")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		t.Fatalf("encode: %v", err)
	}
}

// commandDialogPreview is the New command dialog's body without the pop-up, so
// its footer button row can be inspected for spacing and centring. It is built
// with the same scaffold and button order as showForm, so a capture cannot drift
// from the dialog users see (#503).
func commandDialogPreview() fyne.CanvasObject {
	form := newFormGrid([]formRow{
		{"Name", newScrollFreeEntry()},
		{"Aliases", newScrollFreeEntry()},
		{"Exec", newExecEntry()},
	})

	return container.NewPadded(dialogScaffold(
		"New command",
		container.NewVScroll(form),
		dialogFooter(
			newButton("Cancel", nil, widget.LowImportance, nil),
			newButton("Create", nil, widget.HighImportance, nil),
		),
	))
}

func TestScreenshots(t *testing.T) {
	size := fyne.NewSize(960, 720)

	// Commands: the toolbar rows and an opened card's actions.
	front, _ := commandUI(t, stubExecutor{}, []catalog.Command{
		{Name: "greet", Exec: "echo hi", Icon: "🚀", Group: "deploy", Aliases: []string{"hi"}},
		{Name: "deploy", Exec: "make deploy", Icon: "📦", Group: "deploy"},
	})

	for _, card := range front.commandCards {
		if card.key == "greet" {
			front.toggleCommandCard(card)
			break
		}
	}

	dumpButtons(t, "commands", front.window.Canvas().Content())
	captureShot(t, front.window.Canvas().Content(), "commands", size)

	// History: an opened entry's actions.
	note := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, []string{"a b"})
	hist, _ := historyEntryUI(t, stubExecutor{}, nil, operation.Execute, "greet", note)

	if len(hist.historyCards) > 0 {
		hist.toggleHistoryCard(hist.historyCards[0])
	}

	dumpButtons(t, "history", hist.window.Canvas().Content())
	captureShot(t, hist.window.Canvas().Content(), "history", size)

	// Tools: the card actions and Browse cells.
	tools, _ := commandUI(t, stubExecutor{}, nil)
	tools.tabs.SelectIndex(2)
	captureShot(t, tools.window.Canvas().Content(), "tools", size)

	// About: the update actions, Install revealed and Restore present.
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{},
		&stubUpdater{current: "v0.20.26", previous: "v0.20.25", hasPrevious: true})
	window.SetContent(front.content())
	window.Resize(fyneSize())
	front.applyLatest("v0.20.27", nil)
	front.tabs.SelectIndex(3)
	dumpButtons(t, "about", window.Canvas().Content())
	captureShot(t, window.Canvas().Content(), "about", size)

	// A dialog footer at dialog size.
	dialogPreview := commandDialogPreview()
	dialogWindow := test.NewWindow(dialogPreview)
	dialogWindow.Resize(fyne.NewSize(760, 560))
	dumpButtons(t, "dialog", dialogWindow.Canvas().Content())
	captureShot(t, dialogPreview, "dialog", fyne.NewSize(760, 560))
}
