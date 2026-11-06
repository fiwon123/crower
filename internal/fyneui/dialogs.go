//go:build fyne

package fyneui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
)

// The dialogs are deliberately larger than Fyne's default: a command's exec is
// a shell line, and editing one in a three-line box is the complaint this
// answers.
const (
	dialogWidth      float32 = 760
	dialogHeight     float32 = 560
	editDialogHeight float32 = 680
)

// showCreate asks for the editable fields and adds the command through the
// Editor port. A rule error is shown and the set is left alone.
func (u *ui) showCreate() {
	name := newScrollFreeEntry()
	aliases := newScrollFreeEntry()

	icon, iconWidget := u.iconField("")

	group := newScrollFreeEntry()
	group.SetPlaceHolder("deploy")

	exec := newExecEntry()

	items := []*widget.FormItem{
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Aliases", aliases),
		widget.NewFormItem("Icon", iconWidget),
		widget.NewFormItem("Group", group),
		widget.NewFormItem("Exec", exec),
	}

	// The form scrolls, so a grown exec field or a short window never clips the
	// last row or the buttons.
	form := widget.NewForm(items...)

	showForm("New command", "Create", "Cancel", container.NewVScroll(form), u.window, func(ok bool) {
		if !ok {
			return
		}

		u.applyCreate(name.Text, aliases.Text, icon.value, group.Text, exec.Text)
	}, fyne.NewSize(dialogWidth, dialogHeight))
}

// applyCreate is the form's submit path, separate from the dialog so a test
// can exercise it without building the popup.
func (u *ui) applyCreate(nameText, aliasesText, iconText, groupText, execText string) bool {
	command := catalog.Command{
		Name:    strings.TrimSpace(nameText),
		Aliases: splitAliases(aliasesText),
		Icon:    strings.TrimSpace(iconText),
		Group:   strings.TrimSpace(groupText),
		Exec:    strings.TrimSpace(execText),
	}

	if command.Name == "" || command.Exec == "" {
		showError("Name and exec are required.", u.window)
		return false
	}

	commands, history, err := u.editor.Create(command)
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

	name := newScrollFreeEntry()
	name.SetText(command.Name)

	aliases := newScrollFreeEntry()
	aliases.SetText(strings.Join(command.Aliases, ", "))

	icon, iconWidget := u.iconField(command.Icon)

	group := newScrollFreeEntry()
	group.SetText(command.Group)
	group.SetPlaceHolder("deploy")

	exec := newExecEntry()
	exec.SetText(command.Exec)

	description := newGrowingEntry(3, 8)
	description.SetText(command.Description)
	description.SetPlaceHolder("what the command is for")

	params := newScrollFreeEntry()
	params.SetText(strings.Join(command.Params, ", "))
	params.SetPlaceHolder("env, region")

	cwd := newScrollFreeEntry()
	cwd.SetText(command.Cwd)
	cwd.SetPlaceHolder("~/src/app")

	shell := newScrollFreeEntry()
	shell.SetText(command.Shell)
	shell.SetPlaceHolder("sh, pwsh or cmd")

	timeout := newScrollFreeEntry()
	timeout.SetText(command.Timeout)
	timeout.SetPlaceHolder("120s")

	items := []*widget.FormItem{
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Aliases", aliases),
		widget.NewFormItem("Icon", iconWidget),
		widget.NewFormItem("Group", group),
		widget.NewFormItem("Exec", exec),
		widget.NewFormItem("Description", description),
		widget.NewFormItem("Params", params),
		widget.NewFormItem("Cwd", cwd),
		widget.NewFormItem("Shell", shell),
		widget.NewFormItem("Timeout", timeout),
	}

	form := widget.NewForm(items...)

	showForm("Edit command", "Save", "Cancel", container.NewVScroll(form), u.window, func(ok bool) {
		if !ok {
			return
		}

		u.applyEdit(key, commandForm{
			Name:        name.Text,
			Aliases:     aliases.Text,
			Icon:        icon.value,
			Group:       group.Text,
			Exec:        exec.Text,
			Description: description.Text,
			Params:      params.Text,
			Cwd:         cwd.Text,
			Shell:       shell.Text,
			Timeout:     timeout.Text,
		})
	}, fyne.NewSize(dialogWidth, editDialogHeight))
}

// commandForm is what the edit dialog holds. Every field is always supplied by
// the dialog, so emptying one clears it on the command -- which is what a user
// who deletes the text means.
type commandForm struct {
	Name        string
	Aliases     string
	Icon        string
	Group       string
	Exec        string
	Description string
	Params      string
	Cwd         string
	Shell       string
	Timeout     string
}

// applyEdit is the form's submit path, separate from the dialog so a test can
// exercise it without building the popup.
func (u *ui) applyEdit(key string, form commandForm) bool {
	aliasList := splitAliases(form.Aliases)
	iconText := strings.TrimSpace(form.Icon)
	groupText := strings.TrimSpace(form.Group)

	description := strings.TrimSpace(form.Description)
	paramsList := splitList(form.Params)
	cwd := strings.TrimSpace(form.Cwd)
	shellText := strings.TrimSpace(form.Shell)
	timeout := strings.TrimSpace(form.Timeout)

	changes := catalog.Changes{
		Name:         strings.TrimSpace(form.Name),
		Aliases:      aliasList,
		Exec:         strings.TrimSpace(form.Exec),
		Icon:         iconText,
		Group:        groupText,
		ClearAliases: len(aliasList) == 0,
		ClearIcon:    iconText == "",
		ClearGroup:   groupText == "",

		// The pointer fields say "the dialog supplied this", so an emptied
		// field clears the stored value.
		Description: &description,
		Params:      &paramsList,
		Cwd:         &cwd,
		Shell:       &shellText,
		Timeout:     &timeout,
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

	showConfirm("Delete command", question, u.window, func(ok bool) {
		if !ok {
			return
		}

		u.applyDelete(key)
	})
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

// newExecEntry is the multi-line field a shell line needs. It grows with the
// script, up to twelve rows.
func newExecEntry() *widget.Entry {
	exec := newGrowingEntry(4, 12)
	exec.SetPlaceHolder("echo hello")

	return exec
}

// askParams prompts for the parameters a command declares. The prompt names
// them so a reader knows what the positional arguments mean; an empty answer
// runs the command without parameters.
func (u *ui) askParams(key string, command catalog.Command) {
	u.model.StartRunKey(key, true)

	input := newScrollFreeEntry()
	input.SetPlaceHolder(paramsHint(command))

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

// showDismissible shows content in a non-modal pop-up: a click outside closes
// it. Dialogs that need a decision use showConfirm instead, which stays modal.
func showDismissible(title string, content fyne.CanvasObject, window fyne.Window, size fyne.Size) *widget.PopUp {
	closeButton := widget.NewButton("Close", nil)

	head := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	body := container.NewBorder(
		container.NewVBox(head, widget.NewSeparator()),
		container.NewHBox(layout.NewSpacer(), closeButton),
		nil,
		nil,
		content,
	)

	pop := widget.NewPopUp(container.NewPadded(body), window.Canvas())
	closeButton.OnTapped = pop.Hide
	pop.Resize(size)
	pop.Show()

	return pop
}

// showForm is showDismissible with a confirm/dismiss row. The dismiss button or
// a click outside discards the form; only confirm calls onConfirm.
func showForm(title, confirm, dismiss string, content fyne.CanvasObject, window fyne.Window, onConfirm func(bool), size fyne.Size) *widget.PopUp {
	var pop *widget.PopUp

	confirmButton := widget.NewButton(confirm, func() {
		pop.Hide()
		onConfirm(true)
	})

	dismissButton := widget.NewButton(dismiss, func() {
		if pop != nil {
			pop.Hide()
		}
	})

	head := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	body := container.NewBorder(
		container.NewVBox(head, widget.NewSeparator()),
		container.NewHBox(layout.NewSpacer(), dismissButton, confirmButton),
		nil,
		nil,
		content,
	)

	pop = widget.NewPopUp(container.NewPadded(body), window.Canvas())
	pop.Resize(size)
	pop.Show()

	return pop
}

// showConfirm asks a question. It stays modal so an outside click cannot answer
// it, and its body is a selectable label so the message can be copied.
func showConfirm(title, message string, window fyne.Window, onConfirm func(bool)) {
	dialog.ShowCustomConfirm(title, "Yes", "No", newSelectableText(message), onConfirm, window)
}

func showText(title, text string, window fyne.Window) *widget.PopUp {
	label := newSelectableText(text)
	content := container.NewVScroll(label)

	return showDismissible(title, content, window, fyne.NewSize(dialogWidth, dialogHeight))
}

// newSelectableText is a wrapping label whose text can be selected and copied,
// for output and errors: a user needs to be able to share what happened.
func newSelectableText(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	label.Selectable = true

	return label
}

func showError(message string, window fyne.Window) {
	showText("Error", message, window)
}
