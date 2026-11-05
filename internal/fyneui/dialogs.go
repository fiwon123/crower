//go:build fyne

package fyneui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
)

// The dialogs are deliberately larger than Fyne's default: a command's exec is
// a shell line, and editing one in a three-line box is the complaint this
// answers.
const (
	dialogWidth  float32 = 760
	dialogHeight float32 = 560
)

// showCreate asks for the editable fields and adds the command through the
// Editor port. A rule error is shown and the set is left alone.
func (u *ui) showCreate() {
	name := widget.NewEntry()
	aliases := widget.NewEntry()
	exec := newExecEntry()

	items := []*widget.FormItem{
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Aliases", aliases),
		widget.NewFormItem("Exec", exec),
	}

	form := dialog.NewForm("New command", "Create", "Cancel", items, func(ok bool) {
		if !ok {
			return
		}

		u.applyCreate(name.Text, aliases.Text, exec.Text)
	}, u.window)

	form.Resize(fyne.NewSize(dialogWidth, dialogHeight))
	form.Show()
}

// applyCreate is the form's submit path, separate from the dialog so a test
// can exercise it without building the popup.
func (u *ui) applyCreate(nameText, aliasesText, execText string) bool {
	nameText = strings.TrimSpace(nameText)
	execText = strings.TrimSpace(execText)

	if nameText == "" || execText == "" {
		showError("Name and exec are required.", u.window)
		return false
	}

	commands, history, err := u.editor.Create(nameText, splitAliases(aliasesText), execText)
	if err != nil {
		showError(err.Error(), u.window)
		return false
	}

	u.reload(commands, history)
	return true
}

// showEdit asks for the editable fields of one command and applies the change
// through the Editor port. A rule error is shown and the set is left alone.
func (u *ui) showEditFor(key string) {
	command, found := u.model.Commands.Get(key)
	if !found {
		return
	}

	name := widget.NewEntry()
	name.SetText(command.Name)

	aliases := widget.NewEntry()
	aliases.SetText(strings.Join(command.Aliases, ", "))

	exec := newExecEntry()
	exec.SetText(command.Exec)

	items := []*widget.FormItem{
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Aliases", aliases),
		widget.NewFormItem("Exec", exec),
	}

	// The command's other fields are shown but not editable here: changing them
	// is a schema change with its own phase, and a field that silently did
	// nothing would be worse than a read-only line.
	items = append(items, metadataItems(command)...)

	form := dialog.NewForm("Edit command", "Save", "Cancel", items, func(ok bool) {
		if !ok {
			return
		}

		u.applyEdit(key, name.Text, aliases.Text, exec.Text)
	}, u.window)

	form.Resize(fyne.NewSize(dialogWidth, dialogHeight))
	form.Show()
}

// applyEdit is the form's submit path, separate from the dialog so a test can
// exercise it without building the popup.
func (u *ui) applyEdit(key, nameText, aliasesText, execText string) bool {
	aliasList := splitAliases(aliasesText)
	changes := catalog.Changes{
		Name:         strings.TrimSpace(nameText),
		Aliases:      aliasList,
		Exec:         strings.TrimSpace(execText),
		ClearAliases: len(aliasList) == 0,
	}

	if changes.Name == "" || changes.Exec == "" {
		showError("Name and exec are required.", u.window)
		return false
	}

	commands, history, err := u.editor.Update(key, changes)
	if err != nil {
		showError(err.Error(), u.window)
		return false
	}

	u.reload(commands, history)
	return true
}

// showDelete asks for confirmation and removes the selected command.
func (u *ui) showDeleteFor(key string) {
	if key == "" {
		return
	}

	question := fmt.Sprintf("Delete %q?\n\nIts aliases go with it. You can restore it from the History tab.", key)

	dialog.ShowConfirm("Delete command", question, func(ok bool) {
		if !ok {
			return
		}

		u.applyDelete(key)
	}, u.window)
}

// applyDelete is the confirmation's submit path.
func (u *ui) applyDelete(key string) bool {
	commands, history, err := u.editor.Delete(key)
	if err != nil {
		showError(err.Error(), u.window)
		return false
	}

	u.reload(commands, history)
	return true
}

// newExecEntry is the multi-line field a shell line needs.
func newExecEntry() *widget.Entry {
	exec := widget.NewMultiLineEntry()
	exec.SetMinRowsVisible(6)
	exec.SetPlaceHolder("echo hello")

	return exec
}

// metadataItems renders the fields the dialog cannot change yet, so what the
// command carries is visible while editing it.
func metadataItems(command catalog.Command) []*widget.FormItem {
	items := make([]*widget.FormItem, 0, 5)

	add := func(label, value string) {
		text := widget.NewLabel(value)
		text.Wrapping = fyne.TextWrapWord
		items = append(items, widget.NewFormItem(label, text))
	}

	if command.Description != "" {
		add("About", command.Description)
	}

	if len(command.Params) > 0 {
		add("Params", strings.Join(command.Params, ", "))
	}

	if command.Cwd != "" {
		add("Cwd", command.Cwd)
	}

	if command.Shell != "" {
		add("Shell", command.Shell)
	}

	if command.Timeout != "" {
		add("Timeout", command.Timeout)
	}

	return items
}

func (u *ui) askParams() {
	input := widget.NewEntry()
	input.SetPlaceHolder("space-separated parameters")

	dialog.ShowCustomConfirm("Parameters", "Run", "Cancel", input, func(ok bool) {
		if !ok {
			u.model.State = desktop.ViewList
			return
		}

		u.model.ParamInput = input.Text
		u.model.SubmitParams()
		u.showOutput()
	}, u.window)
}

func (u *ui) showOutput() {
	text := u.model.LastOut
	if u.model.LastErr != nil {
		if text != "" && text[len(text)-1] != '\n' {
			text += "\n"
		}

		text += "error: " + u.model.LastErr.Error()
	}

	if text == "" {
		text = "(no output)"
	}

	title := "Output"
	if u.model.LastCmd != "" {
		title += " — " + u.model.LastCmd
	}

	showText(title, text, u.window)
}

func showText(title, text string, window fyne.Window) {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord

	content := container.NewVScroll(label)

	show := dialog.NewCustom(title, "Close", content, window)
	show.Resize(fyne.NewSize(dialogWidth, dialogHeight))
	show.Show()
}

func showError(message string, window fyne.Window) {
	dialog.ShowError(fmt.Errorf("%s", message), window)
}
