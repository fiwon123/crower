//go:build fyne

package fyneui

import (
	"errors"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

type stubExecutor struct{}

func (stubExecutor) Execute(string, []string) (string, string, error) {
	return "", "", nil
}

type stubEditor struct{}

func (stubEditor) Update(string, catalog.Changes) (catalog.Catalog, timeline.Log, error) {
	return catalog.New(nil, nil), timeline.Log{}, nil
}

func (stubEditor) Create(string, []string, string) (catalog.Catalog, timeline.Log, error) {
	return catalog.New(nil, nil), timeline.Log{}, nil
}

func (stubEditor) Delete(string) (catalog.Catalog, timeline.Log, error) {
	return catalog.New(nil, nil), timeline.Log{}, nil
}

// recordingEditor records the form's calls and returns a canned next state.
type recordingEditor struct {
	updates []catalog.Changes
	creates []string
	deletes []string
	next    catalog.Catalog
	history timeline.Log
	err     error
}

func (e *recordingEditor) Update(_ string, changes catalog.Changes) (catalog.Catalog, timeline.Log, error) {
	e.updates = append(e.updates, changes)

	if e.err != nil {
		return catalog.Catalog{}, timeline.Log{}, e.err
	}

	return e.next, e.history, nil
}

func (e *recordingEditor) Create(name string, _ []string, _ string) (catalog.Catalog, timeline.Log, error) {
	e.creates = append(e.creates, name)

	if e.err != nil {
		return catalog.Catalog{}, timeline.Log{}, e.err
	}

	return e.next, e.history, nil
}

func (e *recordingEditor) Delete(key string) (catalog.Catalog, timeline.Log, error) {
	e.deletes = append(e.deletes, key)

	if e.err != nil {
		return catalog.Catalog{}, timeline.Log{}, e.err
	}

	return e.next, e.history, nil
}

func testModel() *desktop.Model {
	return desktop.NewModel(
		catalog.New(nil, []catalog.Command{{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"}}),
		timeline.Log{},
		stubExecutor{},
	)
}

func testWindow(t *testing.T) fyne.Window {
	t.Helper()

	app := test.NewApp()
	t.Cleanup(app.Quit)

	return test.NewWindow(nil)
}

func listFor(model *desktop.Model) *widget.List {
	return widget.NewList(
		func() int { return len(model.FilteredKeys()) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(widget.ListItemID, fyne.CanvasObject) {},
	)
}

func findButton(object fyne.CanvasObject, label string) *widget.Button {
	if button, ok := object.(*widget.Button); ok && button.Text == label {
		return button
	}

	if container, ok := object.(*fyne.Container); ok {
		for _, child := range container.Objects {
			if button := findButton(child, label); button != nil {
				return button
			}
		}
	}

	return nil
}

// The view builds over the headless model and attaches to a window. The model
// itself is covered by internal/desktop; this is the adapter smoke test.
func TestNewContentBuildsOverTheModel(t *testing.T) {
	window := testWindow(t)

	model := testModel()
	content := NewContent(model, window, stubEditor{})
	if content == nil {
		t.Fatal("NewContent returned nil")
	}

	window.SetContent(content)
	window.Resize(fyne.NewSize(720, 520))

	if window.Canvas().Content() == nil {
		t.Fatal("window has no content")
	}
}

// Every button that asks a question has to put a dialog on the window. The
// dialog content itself is exercised by the apply functions below.
func TestButtonsOpenADialog(t *testing.T) {
	for _, label := range []string{"New", "Edit", "Delete", "History"} {
		t.Run(label, func(t *testing.T) {
			model := testModel()
			window := testWindow(t)

			window.SetContent(NewContent(model, window, &recordingEditor{next: model.Commands}))
			window.Resize(fyne.NewSize(720, 520))

			button := findButton(window.Canvas().Content(), label)
			if button == nil {
				t.Fatalf("no %q button in the window", label)
			}

			test.Tap(button)

			if len(window.Canvas().Overlays().List()) == 0 {
				t.Errorf("%s did not open a dialog", label)
			}
		})
	}
}

func TestApplyEditReloadsTheModel(t *testing.T) {
	model := testModel()
	editor := &recordingEditor{next: catalog.New(nil, []catalog.Command{
		{Name: "ship", Aliases: []string{"s"}, Exec: "echo ship"},
	})}
	window := testWindow(t)

	if !applyEdit(model, editor, "greet", "ship", "s", "echo ship", listFor(model), window) {
		t.Fatal("applyEdit returned false")
	}

	if len(editor.updates) != 1 || editor.updates[0].Name != "ship" || editor.updates[0].Exec != "echo ship" {
		t.Fatalf("updates = %+v, want the edited fields", editor.updates)
	}

	if editor.updates[0].ClearAliases {
		t.Error("ClearAliases = true, want it false for a supplied alias")
	}

	if got := model.Selected(); got != "ship" {
		t.Errorf("Selected() = %q, want the reloaded command", got)
	}
}

func TestApplyEditRejectsEmptyFields(t *testing.T) {
	model := testModel()
	editor := &recordingEditor{next: model.Commands}
	window := testWindow(t)

	if applyEdit(model, editor, "greet", " ", " ", " ", listFor(model), window) {
		t.Fatal("applyEdit accepted an empty name and exec")
	}

	if len(editor.updates) != 0 {
		t.Errorf("updates = %+v, want none", editor.updates)
	}

	if got := model.Selected(); got != "greet" {
		t.Errorf("Selected() = %q, want the set unchanged", got)
	}
}

func TestApplyEditReportsAnEditorError(t *testing.T) {
	model := testModel()
	editor := &recordingEditor{next: model.Commands, err: errors.New("boom")}
	window := testWindow(t)

	if applyEdit(model, editor, "greet", "ship", "s", "echo ship", listFor(model), window) {
		t.Fatal("applyEdit accepted an editor error")
	}

	if got := model.Selected(); got != "greet" {
		t.Errorf("Selected() = %q, want the set unchanged", got)
	}
}

func TestApplyCreateReloadsTheModel(t *testing.T) {
	model := testModel()
	editor := &recordingEditor{next: catalog.New(nil, []catalog.Command{
		{Name: "ship", Exec: "echo ship"},
	})}
	window := testWindow(t)

	if !applyCreate(model, editor, "ship", "", "echo ship", listFor(model), window) {
		t.Fatal("applyCreate returned false")
	}

	if len(editor.creates) != 1 || editor.creates[0] != "ship" {
		t.Fatalf("creates = %v, want ship", editor.creates)
	}

	if got := model.Selected(); got != "ship" {
		t.Errorf("Selected() = %q, want the created command", got)
	}
}

func TestApplyCreateRejectsEmptyFields(t *testing.T) {
	model := testModel()
	editor := &recordingEditor{next: model.Commands}
	window := testWindow(t)

	if applyCreate(model, editor, " ", "", " ", listFor(model), window) {
		t.Fatal("applyCreate accepted an empty name and exec")
	}

	if len(editor.creates) != 0 {
		t.Errorf("creates = %v, want none", editor.creates)
	}
}

func TestApplyDeleteReloadsTheModel(t *testing.T) {
	model := testModel()
	editor := &recordingEditor{next: catalog.New(nil, nil)}
	window := testWindow(t)

	if !applyDelete(model, editor, "greet", listFor(model), window) {
		t.Fatal("applyDelete returned false")
	}

	if len(editor.deletes) != 1 || editor.deletes[0] != "greet" {
		t.Fatalf("deletes = %v, want greet", editor.deletes)
	}

	if got := model.Selected(); got != "" {
		t.Errorf("Selected() = %q, want an empty set", got)
	}
}
