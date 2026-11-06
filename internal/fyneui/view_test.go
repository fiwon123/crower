//go:build fyne

package fyneui

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

type stubExecutor struct{}

func (stubExecutor) Execute(string, []string) (string, string, error) {
	return "", "", nil
}

// recordingExecutor records which command each run asked for.
type recordingExecutor struct {
	keys []string
	args [][]string
}

func (e *recordingExecutor) Execute(key string, params []string) (string, string, error) {
	e.keys = append(e.keys, key)
	e.args = append(e.args, params)

	return "ok", key, nil
}

type stubEditor struct{}

func (stubEditor) Update(string, catalog.Changes) (catalog.Catalog, timeline.Log, error) {
	return catalog.New(nil, nil), timeline.Log{}, nil
}

func (stubEditor) Create(catalog.Command) (catalog.Catalog, timeline.Log, error) {
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

func (e *recordingEditor) Create(command catalog.Command) (catalog.Catalog, timeline.Log, error) {
	e.creates = append(e.creates, command.Name)

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

type stubReverter struct{}

func (stubReverter) Revert(int) (catalog.Catalog, timeline.Log, error) {
	return catalog.New(nil, nil), timeline.Log{}, nil
}

// recordingReverter records the step count the window asked for.
type recordingReverter struct {
	steps   []int
	next    catalog.Catalog
	history timeline.Log
	err     error
}

func (r *recordingReverter) Revert(steps int) (catalog.Catalog, timeline.Log, error) {
	r.steps = append(r.steps, steps)

	if r.err != nil {
		return catalog.Catalog{}, timeline.Log{}, r.err
	}

	return r.next, r.history, nil
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

// newTestUI builds the window's state without the tab container, so a test can
// reach the widgets each tab owns.
func newTestUI(t *testing.T, model *desktop.Model, editor Editor, reverter Reverter) (*ui, fyne.Window) {
	t.Helper()

	return newTestUIWithUpdater(t, model, editor, reverter, &stubUpdater{current: "v0.0.0-test"})
}

func newTestUIWithUpdater(t *testing.T, model *desktop.Model, editor Editor, reverter Reverter, updater Updater) (*ui, fyne.Window) {
	t.Helper()

	window := testWindow(t)
	front := &ui{
		model:       model,
		window:      window,
		editor:      editor,
		reverter:    reverter,
		updater:     updater,
		tools:       &stubTools{},
		groups:      &stubGroups{},
		setups:      &stubSetups{current: "default", names: []string{"default"}},
		activeSetup: "default",
		groupIcons:  map[string]string{},
	}

	return front, window
}

// The window is tab based: Commands, History, Tools and About today, with the
// setup switcher adding its own tab.
func TestNewContentBuildsTheTabs(t *testing.T) {
	window := testWindow(t)
	model := testModel()
	updater := &stubUpdater{current: "v0.0.0-test"}

	content := NewContent(model, window, Deps{
		Exec:     stubExecutor{},
		Editor:   stubEditor{},
		Reverter: stubReverter{},
		Updater:  updater,
		Tools:    &stubTools{},
		Groups:   &stubGroups{},
		Setups:   &stubSetups{current: "default", names: []string{"default"}},
	})

	tabs, ok := content.(*container.AppTabs)
	if !ok {
		t.Fatalf("content = %T, want tabs", content)
	}

	want := []string{"Commands", "History", "Tools", "About"}
	if len(tabs.Items) != len(want) {
		t.Fatalf("got %d tabs, want %d", len(tabs.Items), len(want))
	}

	for i, text := range want {
		if tabs.Items[i].Text != text {
			t.Errorf("tab %d = %q, want %q", i, tabs.Items[i].Text, text)
		}
	}

	window.SetContent(content)
	window.Resize(fyne.NewSize(960, 640))

	if window.Canvas().Content() == nil {
		t.Fatal("window has no content")
	}
}

// labelsIn collects every label text the canvas can reach, walking the tree the
// way Fyne paints it: a container exposes its objects, and a widget exposes its
// renderer's objects. Anything else is opaque — including a type that embeds a
// container, which is the shape that made the rows paint nothing (#406).
func labelsIn(object fyne.CanvasObject) []string {
	switch typed := object.(type) {
	case *widget.Label:
		return []string{typed.Text}
	case *fyne.Container:
		found := make([]string, 0, len(typed.Objects))

		for _, child := range typed.Objects {
			found = append(found, labelsIn(child)...)
		}

		return found
	case fyne.Widget:
		renderer := test.WidgetRenderer(typed)
		if renderer == nil {
			return nil
		}

		found := make([]string, 0, len(renderer.Objects()))

		for _, child := range renderer.Objects() {
			found = append(found, labelsIn(child)...)
		}

		return found
	}

	return nil
}

// labelPosition finds a label by text and returns its absolute position.
func labelPosition(object fyne.CanvasObject, text string) (fyne.Position, bool) {
	switch typed := object.(type) {
	case *widget.Label:
		if typed.Text == text {
			return fyne.CurrentApp().Driver().AbsolutePositionForObject(typed), true
		}
	case *fyne.Container:
		for _, child := range typed.Objects {
			if pos, found := labelPosition(child, text); found {
				return pos, true
			}
		}
	case fyne.Widget:
		renderer := test.WidgetRenderer(typed)
		if renderer == nil {
			return fyne.Position{}, false
		}

		for _, child := range renderer.Objects() {
			if pos, found := labelPosition(child, text); found {
				return pos, true
			}
		}
	}

	return fyne.Position{}, false
}

// commandUI lays out the Commands tab over the given commands.
func commandUI(t *testing.T, exec desktop.Executor, commands []catalog.Command) (*ui, fyne.Window) {
	t.Helper()

	model := desktop.NewModel(catalog.New(nil, commands), timeline.Log{}, exec)
	front, window := newTestUI(t, model, stubEditor{}, stubReverter{})
	window.SetContent(front.content())
	window.Resize(fyneSize())

	return front, window
}

// openCommandCard opens the card for key.
func openCommandCard(t *testing.T, front *ui, key string) *commandCard {
	t.Helper()

	for _, card := range front.commandCards {
		if card.key == key {
			front.toggleCommandCard(card)
			return card
		}
	}

	t.Fatalf("command %q has no card", key)
	return nil
}

// The Commands tab keeps its controls and cards off the window edge.
func TestCommandsHaveALeftMargin(t *testing.T) {
	front, window := newTestUI(t, testModel(), stubEditor{}, stubReverter{})

	window.SetContent(front.content())
	window.Resize(fyneSize())

	pos, found := labelPosition(window.Canvas().Content(), "Setup")
	if !found {
		t.Fatal("the Setup label is not in the rendered tree")
	}

	if pos.X < theme.Padding()*2 {
		t.Errorf("label x = %.1f, want at least %.1f from the window edge", pos.X, theme.Padding()*2)
	}
}

// hasButtonText reports whether a button with exactly this text is reachable.
func hasButtonText(object fyne.CanvasObject, text string) bool {
	switch typed := object.(type) {
	case *widget.Button:
		return typed.Text == text
	case *fyne.Container:
		for _, child := range typed.Objects {
			if hasButtonText(child, text) {
				return true
			}
		}
	case fyne.Widget:
		renderer := test.WidgetRenderer(typed)
		if renderer == nil {
			return false
		}

		for _, child := range renderer.Objects() {
			if hasButtonText(child, text) {
				return true
			}
		}
	}

	return false
}

// The setup and command controls say what they act on.
func TestSetupButtonsAreNamed(t *testing.T) {
	front, window := newTestUI(t, testModel(), stubEditor{}, stubReverter{})

	window.SetContent(front.content())
	window.Resize(fyneSize())

	for _, want := range []string{"New command", "New setup…", "Delete setup…", "Groups…"} {
		if !hasButtonText(window.Canvas().Content(), want) {
			t.Errorf("no %q button in the Commands tab", want)
		}
	}

	if hasButtonText(window.Canvas().Content(), "Details") {
		t.Error("the global Details toggle is still there")
	}
}

// The canvas reaches an opened card's text. A card whose body is a container
// subtype would paint nothing (#406).
func TestCommandCardIsPaintable(t *testing.T) {
	front, window := commandUI(t, stubExecutor{}, []catalog.Command{{Name: "greet", Exec: "echo hi"}})

	openCommandCard(t, front, "greet")

	if labels := labelsIn(window.Canvas().Content()); !slices.Contains(labels, "echo hi") {
		t.Fatalf("the open card does not expose its exec; labels = %q", labels)
	}
}

// A card's Run runs that card's command.
func TestCommandCardRunsItsOwnCommand(t *testing.T) {
	exec := &recordingExecutor{}
	front, _ := commandUI(t, exec, []catalog.Command{
		{Name: "greet", Exec: "echo hi"},
		{Name: "ship", Exec: "echo ship"},
	})

	first := openCommandCard(t, front, "greet")
	if first.run.OnTapped == nil {
		t.Fatal("a card's run button has no action")
	}

	first.run.OnTapped()

	second := openCommandCard(t, front, "ship")
	second.run.OnTapped()

	if !slices.Equal(exec.keys, []string{"greet", "ship"}) {
		t.Fatalf("executed %v, want [greet ship]", exec.keys)
	}
}

// The card's Run takes its parameters from the entry; an empty entry means
// "run without parameters".
func TestCommandCardRunUsesTheParamsEntry(t *testing.T) {
	exec := &recordingExecutor{}
	front, _ := commandUI(t, exec, []catalog.Command{{Name: "greet", Exec: "echo hi"}})

	card := openCommandCard(t, front, "greet")
	card.run.OnTapped()

	if len(exec.args) != 1 || len(exec.args[0]) != 0 {
		t.Fatalf("args = %v, want one run with no parameters", exec.args)
	}

	card.params.SetText("a  b")
	card.run.OnTapped()

	if len(exec.args) != 2 || !slices.Equal(exec.args[1], []string{"a", "b"}) {
		t.Fatalf("args = %v, want [a b]", exec.args)
	}
}

// A collapsed card carries no actions; opening it shows the fields, the params
// entry and Run/Edit/Delete.
func TestCommandCardShowsTheDetail(t *testing.T) {
	front, _ := commandUI(t, stubExecutor{}, []catalog.Command{
		{Name: "deploy", Exec: "make deploy", Params: []string{"env"}},
	})

	if front.commandOpen != nil {
		t.Fatal("a card is open before one was tapped")
	}

	card := openCommandCard(t, front, "deploy")

	if !card.isOpen() || !card.body.Visible() {
		t.Fatal("opening the card did not show its body")
	}

	labels := labelsIn(card.body)
	for _, want := range []string{"Name", "Exec", "make deploy", "Params", "env"} {
		if !slices.Contains(labels, want) {
			t.Errorf("fields %q do not contain %q", labels, want)
		}
	}

	if card.run.Text != "Run" || card.edit.Text != "Edit" || card.remove.Text != "Delete" {
		t.Errorf("actions = %q/%q/%q, want Run/Edit/Delete", card.run.Text, card.edit.Text, card.remove.Text)
	}

	if !strings.Contains(card.params.PlaceHolder, "env") {
		t.Errorf("placeholder = %q, want the declared parameter names", card.params.PlaceHolder)
	}
}

// The empty state distinguishes "no commands" from "nothing matches", and a
// filter rebuilds the cards and drops the open one.
func TestCommandCardsFollowTheFilter(t *testing.T) {
	front, _ := commandUI(t, stubExecutor{}, nil)

	if !front.commandEmpty.Visible() || !strings.Contains(front.commandEmpty.Text, "No commands yet") {
		t.Errorf("empty state = %q, want the no-commands message", front.commandEmpty.Text)
	}

	front, _ = commandUI(t, stubExecutor{}, []catalog.Command{
		{Name: "greet", Exec: "echo hi"},
		{Name: "ship", Exec: "echo ship"},
	})

	openCommandCard(t, front, "greet")

	front.filterEntry.SetText("shi")

	if len(front.commandCards) != 1 || front.commandCards[0].key != "ship" {
		t.Fatalf("cards = %v, want only ship", front.commandCards)
	}

	if front.commandOpen != nil {
		t.Error("a card stayed open after the filter changed")
	}

	front.filterEntry.SetText("")

	if len(front.commandCards) != 2 {
		t.Errorf("cards = %d, want both commands", len(front.commandCards))
	}

	front.filterEntry.SetText("zzz")

	if !front.commandEmpty.Visible() || !strings.Contains(front.commandEmpty.Text, "match") {
		t.Errorf("empty state = %q, want the no-match message", front.commandEmpty.Text)
	}
}

// Opening another card closes the first.
func TestOpeningAnotherCommandCardClosesTheFirst(t *testing.T) {
	front, _ := commandUI(t, stubExecutor{}, []catalog.Command{
		{Name: "greet", Exec: "echo hi"},
		{Name: "ship", Exec: "echo ship"},
	})

	first := openCommandCard(t, front, "greet")
	second := openCommandCard(t, front, "ship")

	if first.isOpen() {
		t.Error("the first card stayed open")
	}

	if !second.isOpen() {
		t.Error("the second card did not open")
	}
}

// Opening a card gives its body a height, so the page can lay it out.
func TestCommandCardSizesTheBody(t *testing.T) {
	front, _ := commandUI(t, stubExecutor{}, []catalog.Command{{Name: "greet", Exec: "echo hi"}})

	card := openCommandCard(t, front, "greet")

	if size := card.body.Size(); size.Height == 0 {
		t.Fatalf("body size = %v, want the page to lay it out", size)
	}
}

// Delete never happens on the first click.
func TestCommandCardDeleteAsksFirst(t *testing.T) {
	front, window := commandUI(t, stubExecutor{}, []catalog.Command{{Name: "greet", Exec: "echo hi"}})

	card := openCommandCard(t, front, "greet")
	if card.remove.OnTapped == nil {
		t.Fatal("the delete button has no action")
	}

	card.remove.OnTapped()

	if len(window.Canvas().Overlays().List()) == 0 {
		t.Error("delete did not ask for confirmation")
	}
}

// Edit opens a dialog rather than changing the command in place.
func TestCommandCardEditOpensTheDialog(t *testing.T) {
	front, window := commandUI(t, stubExecutor{}, []catalog.Command{{Name: "greet", Exec: "echo hi"}})

	card := openCommandCard(t, front, "greet")
	if card.edit.OnTapped == nil {
		t.Fatal("the edit button has no action")
	}

	card.edit.OnTapped()

	if len(window.Canvas().Overlays().List()) == 0 {
		t.Error("edit did not open a dialog")
	}
}

// The detail grid carries the metadata that does not fit on a row, and omits
// what the command does not set.
func TestCommandDetailFields(t *testing.T) {
	fields := commandDetailFields(catalog.Command{
		Name:        "deploy",
		Exec:        "make deploy",
		Aliases:     []string{"d", "ship"},
		Group:       "work",
		Description: "ship to production",
		Params:      []string{"env", "region"},
		Cwd:         "~/src/app",
		Shell:       "sh",
		Timeout:     "120s",
		Steps:       []string{"make build", "make deploy"},
	})

	want := []detailField{
		{Label: "Name", Value: "deploy"},
		{Label: "Exec", Value: "make deploy"},
		{Label: "Aliases", Value: "d, ship"},
		{Label: "Group", Value: "work"},
		{Label: "Description", Value: "ship to production"},
		{Label: "Params", Value: "env, region"},
		{Label: "Cwd", Value: "~/src/app"},
		{Label: "Shell", Value: "sh"},
		{Label: "Timeout", Value: "120s"},
		{Label: "Steps", Value: "1. make build\n2. make deploy"},
	}

	if !slices.Equal(fields, want) {
		t.Errorf("commandDetailFields() = %v, want %v", fields, want)
	}

	if got := commandDetailFields(catalog.Command{Name: "plain", Exec: "echo hi"}); !slices.Equal(got, []detailField{
		{Label: "Name", Value: "plain"},
		{Label: "Exec", Value: "echo hi"},
	}) {
		t.Errorf("plain command = %v, want Name and Exec", got)
	}
}

func TestApplyEditReloadsTheModel(t *testing.T) {
	front, _ := newTestUI(t, testModel(), &recordingEditor{next: catalog.New(nil, []catalog.Command{
		{Name: "ship", Aliases: []string{"s"}, Exec: "echo ship"},
	})}, stubReverter{})

	front.window.SetContent(widget.NewLabel(""))

	if !front.applyEdit("greet", commandForm{
		Name:        "ship",
		Aliases:     "s",
		Icon:        "🚀",
		Group:       "deploy",
		Exec:        "echo ship",
		Description: "ship it",
		Params:      "env, region",
		Cwd:         "/srv/app",
		Shell:       "sh",
		Timeout:     "30s",
	}) {
		t.Fatal("applyEdit returned false")
	}

	editor := front.editor.(*recordingEditor)

	if len(editor.updates) != 1 || editor.updates[0].Name != "ship" || editor.updates[0].Exec != "echo ship" {
		t.Fatalf("updates = %+v, want the edited fields", editor.updates)
	}

	if editor.updates[0].ClearAliases {
		t.Error("ClearAliases = true, want it false for a supplied alias")
	}

	if editor.updates[0].Icon != "🚀" || editor.updates[0].Group != "deploy" {
		t.Errorf("changes = %+v, want the icon and group carried through", editor.updates[0])
	}

	if editor.updates[0].ClearIcon || editor.updates[0].ClearGroup {
		t.Error("a supplied icon or group was marked for clearing")
	}

	changes := editor.updates[0]

	if changes.Description == nil || *changes.Description != "ship it" {
		t.Errorf("Description = %v, want it carried through", changes.Description)
	}

	if changes.Params == nil || !slices.Equal(*changes.Params, []string{"env", "region"}) {
		t.Errorf("Params = %v, want [env region]", changes.Params)
	}

	if changes.Cwd == nil || *changes.Cwd != "/srv/app" {
		t.Errorf("Cwd = %v, want /srv/app", changes.Cwd)
	}

	if changes.Shell == nil || *changes.Shell != "sh" {
		t.Errorf("Shell = %v, want sh", changes.Shell)
	}

	if changes.Timeout == nil || *changes.Timeout != "30s" {
		t.Errorf("Timeout = %v, want 30s", changes.Timeout)
	}

	if got := front.model.Selected(); got != "ship" {
		t.Errorf("Selected() = %q, want the reloaded command", got)
	}
}

func TestApplyEditRejectsEmptyFields(t *testing.T) {
	editor := &recordingEditor{next: testModel().Commands}
	front, _ := newTestUI(t, testModel(), editor, stubReverter{})
	front.window.SetContent(widget.NewLabel(""))

	if front.applyEdit("greet", commandForm{Name: " ", Aliases: " ", Icon: " ", Group: " ", Exec: " "}) {
		t.Fatal("applyEdit accepted an empty name and exec")
	}

	if len(editor.updates) != 0 {
		t.Errorf("updates = %+v, want none", editor.updates)
	}

	if got := front.model.Selected(); got != "greet" {
		t.Errorf("Selected() = %q, want the set unchanged", got)
	}
}

func TestApplyEditReportsAnEditorError(t *testing.T) {
	editor := &recordingEditor{next: testModel().Commands, err: errors.New("boom")}
	front, _ := newTestUI(t, testModel(), editor, stubReverter{})
	front.window.SetContent(widget.NewLabel(""))

	if front.applyEdit("greet", commandForm{Name: "ship", Aliases: "s", Icon: "🚀", Group: "deploy", Exec: "echo ship"}) {
		t.Fatal("applyEdit accepted an editor error")
	}

	if got := front.model.Selected(); got != "greet" {
		t.Errorf("Selected() = %q, want the set unchanged", got)
	}
}

func TestApplyCreateReloadsTheModel(t *testing.T) {
	editor := &recordingEditor{next: catalog.New(nil, []catalog.Command{
		{Name: "ship", Exec: "echo ship"},
	})}
	front, _ := newTestUI(t, testModel(), editor, stubReverter{})
	front.window.SetContent(widget.NewLabel(""))

	if !front.applyCreate("ship", "", "🚀", "deploy", "echo ship") {
		t.Fatal("applyCreate returned false")
	}

	if len(editor.creates) != 1 || editor.creates[0] != "ship" {
		t.Fatalf("creates = %v, want ship", editor.creates)
	}

	if got := front.model.Selected(); got != "ship" {
		t.Errorf("Selected() = %q, want the created command", got)
	}
}

func TestApplyCreateRejectsEmptyFields(t *testing.T) {
	editor := &recordingEditor{next: testModel().Commands}
	front, _ := newTestUI(t, testModel(), editor, stubReverter{})
	front.window.SetContent(widget.NewLabel(""))

	if front.applyCreate(" ", "", "", "", " ") {
		t.Fatal("applyCreate accepted an empty name and exec")
	}

	if len(editor.creates) != 0 {
		t.Errorf("creates = %v, want none", editor.creates)
	}
}

func TestApplyDeleteReloadsTheModel(t *testing.T) {
	editor := &recordingEditor{next: catalog.New(nil, nil)}
	front, _ := newTestUI(t, testModel(), editor, stubReverter{})
	front.window.SetContent(widget.NewLabel(""))

	if !front.applyDelete("greet") {
		t.Fatal("applyDelete returned false")
	}

	if len(editor.deletes) != 1 || editor.deletes[0] != "greet" {
		t.Fatalf("deletes = %v, want greet", editor.deletes)
	}

	if got := front.model.Selected(); got != "" {
		t.Errorf("Selected() = %q, want an empty set", got)
	}
}

func TestApplyRevertReloadsTheModel(t *testing.T) {
	reverter := &recordingReverter{next: catalog.New(nil, []catalog.Command{
		{Name: "greet", Exec: "echo hi"},
	})}
	front, _ := newTestUI(t, testModel(), stubEditor{}, reverter)
	front.window.SetContent(widget.NewLabel(""))

	if !front.applyRevert(2) {
		t.Fatal("applyRevert returned false")
	}

	if len(reverter.steps) != 1 || reverter.steps[0] != 2 {
		t.Fatalf("steps = %v, want [2]", reverter.steps)
	}

	if got := front.model.Selected(); got != "greet" {
		t.Errorf("Selected() = %q, want the restored command", got)
	}
}

func TestApplyRevertReportsAnError(t *testing.T) {
	reverter := &recordingReverter{err: errors.New("boom")}
	front, _ := newTestUI(t, testModel(), stubEditor{}, reverter)
	front.window.SetContent(widget.NewLabel(""))

	if front.applyRevert(1) {
		t.Fatal("applyRevert accepted a reverter error")
	}
}

// A row shows the icon when one is set, and the bare name when it is not.
func TestCommandLabelShowsTheIcon(t *testing.T) {
	if got := commandLabel(catalog.Command{Name: "deploy", Icon: "🚀"}); got != "🚀 deploy" {
		t.Errorf("commandLabel = %q, want the icon before the name", got)
	}

	if got := commandLabel(catalog.Command{Name: "deploy"}); got != "deploy" {
		t.Errorf("commandLabel = %q, want the bare name", got)
	}
}
