//go:build fyne

// Package fyneui is the Fyne front end. It renders the headless desktop model
// and owns no rules of its own; the default binary does not build it, so the
// cgo-free release is unchanged.
package fyneui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// Editor applies an edit, a creation or a deletion and returns the refreshed
// set and history.
type Editor interface {
	Update(key string, changes catalog.Changes) (catalog.Catalog, timeline.Log, error)
	Create(name string, aliases []string, exec string) (catalog.Catalog, timeline.Log, error)
	Delete(key string) (catalog.Catalog, timeline.Log, error)
}

// Reverter restores the command set an older history entry recorded and
// returns the refreshed set and history.
//
// It is declared here, at the point of use: the window confirms the restore
// itself, so it needs the state change and not the use case's prompt.
type Reverter interface {
	Revert(steps int) (catalog.Catalog, timeline.Log, error)
}

// Deps is what the window needs from the command layer. Every field is a port
// declared in this package and implemented in cmd, so the UI depends on
// behaviour, not on internal/session.
type Deps struct {
	Exec     desktop.Executor
	Editor   Editor
	Reverter Reverter
	Updater  Updater
	Tools    desktop.Tools
}

// Run opens the Fyne window and blocks until it closes.
func Run(commands catalog.Catalog, history timeline.Log, deps Deps) {
	model := desktop.NewModel(commands, history, deps.Exec)

	a := app.New()
	a.SetIcon(logoResource())

	window := a.NewWindow("crower")
	front := newUI(model, window, deps)
	window.SetContent(front.content())
	window.Resize(fyne.NewSize(960, 640))

	// The lookup runs after the content exists and never blocks the window:
	// opening crower is not a reason to wait for GitHub.
	front.checkForUpdates()

	window.ShowAndRun()
}

// ui is the widget set the tabs share. Reloading after an edit has to refresh
// both the command list and the history table, so both live here instead of
// being captured by one tab.
type ui struct {
	model    *desktop.Model
	window   fyne.Window
	editor   Editor
	reverter Reverter
	updater  Updater
	tools    desktop.Tools

	list          *widget.List
	detail        *widget.Label
	table         *widget.Table
	historyDetail *widget.Label
	revert        *widget.Button

	// historyRows is the table's snapshot: rebuild it on reload or tab open
	// rather than recomputing the whole log for every cell.
	historyRows []desktop.HistoryRow
	historyRow  int

	// tools tab
	toolResult *widget.Label

	// about tab
	updateStatus  *widget.Label
	checkButton   *widget.Button
	installButton *widget.Button
	restoreButton *widget.Button
	latest        string
}

func newUI(model *desktop.Model, window fyne.Window, deps Deps) *ui {
	return &ui{
		model:      model,
		window:     window,
		editor:     deps.Editor,
		reverter:   deps.Reverter,
		updater:    deps.Updater,
		tools:      deps.Tools,
		historyRow: -1,
	}
}

// NewContent builds the window content over the headless model.
func NewContent(model *desktop.Model, window fyne.Window, deps Deps) fyne.CanvasObject {
	return newUI(model, window, deps).content()
}

// content builds the tabs. It is network-free: only Run starts the update
// lookup, so a test can build the window without reaching GitHub.
func (u *ui) content() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Commands", u.commandsTab()),
		container.NewTabItem("History", u.historyTab()),
		container.NewTabItem("Tools", u.toolsTab()),
		container.NewTabItem("About", u.aboutTab()),
	)

	// The table is built from the model, so a command edit changes it. Refresh
	// when the tab is opened instead of refreshing a hidden table on every
	// keystroke.
	tabs.OnSelected = func(item *container.TabItem) {
		if item.Text == "History" {
			u.refreshHistory()
		}
	}

	return tabs
}

// commandsTab is the command list with its filter and per-row actions.
func (u *ui) commandsTab() fyne.CanvasObject {
	u.detail = widget.NewLabel("")
	u.detail.Wrapping = fyne.TextWrapWord

	u.list = widget.NewList(
		func() int { return len(u.model.FilteredKeys()) },
		func() fyne.CanvasObject { return newCommandRow() },
		func(id widget.ListItemID, object fyne.CanvasObject) { u.updateRow(id, object) },
	)

	u.list.OnSelected = func(id widget.ListItemID) { u.selectRow(id) }
	u.list.OnUnselected = func(widget.ListItemID) { u.detail.SetText("") }

	filter := widget.NewEntry()
	filter.SetPlaceHolder("Filter commands")
	filter.OnChanged = func(string) {
		u.model.Filter = filter.Text
		u.model.Index = 0
		u.list.Refresh()
		u.list.UnselectAll()
		u.detail.SetText("")
	}

	toolbar := container.NewHBox(
		widget.NewButtonWithIcon("New", theme.ContentAddIcon(), u.showCreate),
		widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), func() { u.runSelected(false) }),
		widget.NewButtonWithIcon("Run with params…", theme.MediaFastForwardIcon(), func() { u.runSelected(true) }),
		widget.NewButton("Quit", u.window.Close),
	)

	return container.NewBorder(container.NewVBox(filter, u.detail), toolbar, nil, nil, u.list)
}

// commandRow is one list row: the command name on the left, its actions on the
// right. widget.List recycles templates, so updateRow re-points every button at
// the row's own key on each update.
type commandRow struct {
	*fyne.Container

	name   *widget.Label
	run    *widget.Button
	params *widget.Button
	edit   *widget.Button
	remove *widget.Button
}

func newCommandRow() *commandRow {
	row := &commandRow{
		name:   widget.NewLabel(""),
		run:    widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
		params: widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), nil),
		edit:   widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil),
		remove: widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
	}

	row.run.Importance = widget.LowImportance
	row.params.Importance = widget.LowImportance
	row.edit.Importance = widget.LowImportance
	row.remove.Importance = widget.DangerImportance

	row.Container = container.NewBorder(
		nil,
		nil,
		nil,
		container.NewHBox(row.run, row.params, row.edit, row.remove),
		row.name,
	)

	return row
}

func (u *ui) updateRow(id widget.ListItemID, object fyne.CanvasObject) {
	row, ok := object.(*commandRow)
	if !ok {
		return
	}

	keys := u.model.FilteredKeys()
	if int(id) >= len(keys) {
		row.name.SetText("")
		return
	}

	key := keys[id]
	row.name.SetText(key)

	row.run.OnTapped = func() { u.runKey(key, false) }
	row.params.OnTapped = func() { u.runKey(key, true) }
	row.edit.OnTapped = func() { u.showEditFor(key) }
	row.remove.OnTapped = func() { u.showDeleteFor(key) }
}

func (u *ui) selectRow(id widget.ListItemID) {
	u.model.Index = int(id)

	selected := u.model.Selected()
	if selected == "" {
		u.detail.SetText("")
		return
	}

	command, found := u.model.Commands.Get(selected)
	if !found {
		u.detail.SetText(selected)
		return
	}

	u.detail.SetText(describe(command))
}

// describe is the detail line: what the command runs, plus the metadata that
// does not fit on a row.
func describe(command catalog.Command) string {
	lines := []string{command.Name, command.Exec}

	if len(command.Aliases) > 0 {
		lines = append(lines, "aliases: "+strings.Join(command.Aliases, ", "))
	}

	if command.Description != "" {
		lines = append(lines, "about: "+command.Description)
	}

	if len(command.Params) > 0 {
		lines = append(lines, "params: "+strings.Join(command.Params, ", "))
	}

	if command.Cwd != "" {
		lines = append(lines, "cwd: "+command.Cwd)
	}

	if command.Shell != "" {
		lines = append(lines, "shell: "+command.Shell)
	}

	if command.Timeout != "" {
		lines = append(lines, "timeout: "+command.Timeout)
	}

	if len(command.Steps) > 0 {
		lines = append(lines, fmt.Sprintf("%d steps", len(command.Steps)))
	}

	return strings.Join(lines, "\n")
}

func (u *ui) runSelected(withParams bool) {
	u.runKey(u.model.Selected(), withParams)
}

// runKey runs one command: the row's key when a row button asked, the
// selection when the toolbar did.
func (u *ui) runKey(key string, withParams bool) {
	if key == "" {
		return
	}

	u.model.StartRunKey(key, withParams)

	if withParams {
		u.askParams()
		return
	}

	u.showOutput()
}

// reload replaces the model state after an edit and refreshes every view built
// from it. The widgets are optional so a test can exercise the submit paths
// without a window.
func (u *ui) reload(commands catalog.Catalog, history timeline.Log) {
	u.model.Reload(commands, history)

	if u.list != nil {
		u.list.Refresh()
		u.list.UnselectAll()
	}

	if u.detail != nil {
		u.detail.SetText("")
	}

	u.refreshHistory()
}

func splitAliases(text string) []string {
	aliases := make([]string, 0)

	for _, part := range strings.Split(text, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			aliases = append(aliases, part)
		}
	}

	return aliases
}
