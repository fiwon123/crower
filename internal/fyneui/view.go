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
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// Editor applies an edit and returns the refreshed set and history.
type Editor interface {
	Update(key string, changes catalog.Changes) (catalog.Catalog, timeline.Log, error)
}

// Run opens the Fyne window and blocks until it closes.
func Run(commands catalog.Catalog, history timeline.Log, exec desktop.Executor, editor Editor) {
	model := desktop.NewModel(commands, history, exec)

	window := app.New().NewWindow("crower fyne")
	window.SetContent(NewContent(model, window, editor))
	window.Resize(fyne.NewSize(720, 520))
	window.ShowAndRun()
}

// NewContent builds the window content over the headless model.
func NewContent(model *desktop.Model, window fyne.Window, editor Editor) fyne.CanvasObject {
	keys := widget.NewLabel("")
	keys.Wrapping = fyne.TextWrapWord

	list := widget.NewList(
		func() int { return len(model.FilteredKeys()) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, object fyne.CanvasObject) {
			filtered := model.FilteredKeys()
			if int(id) >= len(filtered) {
				return
			}

			object.(*widget.Label).SetText(filtered[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		model.Index = int(id)
		selected := model.Selected()
		if selected == "" {
			keys.SetText("")
			return
		}

		command, found := model.Commands.Get(selected)
		if !found {
			keys.SetText(selected)
			return
		}

		keys.SetText(fmt.Sprintf("%s\n%s", command.Name, command.Exec))
	}

	filter := widget.NewEntry()
	filter.SetPlaceHolder("Filter commands")
	filter.OnChanged = func(string) {
		model.Filter = filter.Text
		model.Index = 0
		list.Refresh()
		list.UnselectAll()
	}

	run := widget.NewButton("Run", func() {
		model.StartRun(false)
		showOutput(model, window)
	})

	params := widget.NewButton("Run with params…", func() {
		model.StartRun(true)
		askParams(model, window)
	})

	edit := widget.NewButton("Edit", func() {
		showEdit(model, editor, window, list)
	})

	history := widget.NewButton("History", func() {
		showText("History", model.HistoryText(), window)
	})

	quit := widget.NewButton("Quit", func() { window.Close() })

	buttons := container.NewHBox(run, params, edit, history, quit)

	return container.NewBorder(
		container.NewVBox(filter, keys),
		buttons,
		nil,
		nil,
		list,
	)
}

// showEdit asks for the three editable fields and applies the change through
// the Editor port. A rule error is shown and the set is left alone.
func showEdit(model *desktop.Model, editor Editor, window fyne.Window, list *widget.List) {
	key := model.Selected()
	if key == "" {
		return
	}

	command, found := model.Commands.Get(key)
	if !found {
		return
	}

	name := widget.NewEntry()
	name.SetText(command.Name)

	aliases := widget.NewEntry()
	aliases.SetText(strings.Join(command.Aliases, ", "))

	exec := widget.NewEntry()
	exec.SetText(command.Exec)

	items := []*widget.FormItem{
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Aliases", aliases),
		widget.NewFormItem("Exec", exec),
	}

	dialog.NewForm("Edit command", "Save", "Cancel", items, func(ok bool) {
		if !ok {
			return
		}

		aliasList := splitAliases(aliases.Text)
		changes := catalog.Changes{
			Name:         strings.TrimSpace(name.Text),
			Aliases:      aliasList,
			Exec:         strings.TrimSpace(exec.Text),
			ClearAliases: len(aliasList) == 0,
		}

		if changes.Name == "" || changes.Exec == "" {
			showError("Name and exec are required.", window)
			return
		}

		commands, history, err := editor.Update(key, changes)
		if err != nil {
			showError(err.Error(), window)
			return
		}

		model.Reload(commands, history)
		list.Refresh()
		list.UnselectAll()
	}, window).Show()
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

func showError(message string, window fyne.Window) {
	dialog.ShowError(fmt.Errorf("%s", message), window)
}

func askParams(model *desktop.Model, window fyne.Window) {
	input := widget.NewEntry()
	input.SetPlaceHolder("space-separated parameters")

	dialog.ShowCustomConfirm("Parameters", "Run", "Cancel", input, func(ok bool) {
		if !ok {
			model.State = desktop.ViewList
			return
		}

		model.ParamInput = input.Text
		model.SubmitParams()
		showOutput(model, window)
	}, window)
}

func showOutput(model *desktop.Model, window fyne.Window) {
	text := model.LastOut
	if model.LastErr != nil {
		if text != "" && text[len(text)-1] != '\n' {
			text += "\n"
		}

		text += "error: " + model.LastErr.Error()
	}

	if text == "" {
		text = "(no output)"
	}

	title := "Output"
	if model.LastCmd != "" {
		title += " — " + model.LastCmd
	}

	showText(title, text, window)
}

func showText(title, text string, window fyne.Window) {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord

	content := container.NewVScroll(label)

	dialog.NewCustom(title, "Close", content, window).Show()
}
