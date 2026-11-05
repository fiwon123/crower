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
	dialogWidth      float32 = 760
	dialogHeight     float32 = 560
	editDialogHeight float32 = 680
)

// showCreate asks for the editable fields and adds the command through the
// Editor port. A rule error is shown and the set is left alone.
func (u *ui) showCreate() {
	name := newScrollFreeEntry()
	aliases := newScrollFreeEntry()

	icon := newScrollFreeEntry()
	icon.SetPlaceHolder("🚀")

	group := newScrollFreeEntry()
	group.SetPlaceHolder("deploy")

	exec := newExecEntry()

	items := []*widget.FormItem{
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Aliases", aliases),
		widget.NewFormItem("Icon", u.iconField(icon)),
		widget.NewFormItem("Group", group),
		widget.NewFormItem("Exec", exec),
	}

	form := dialog.NewForm("New command", "Create", "Cancel", items, func(ok bool) {
		if !ok {
			return
		}

		u.applyCreate(name.Text, aliases.Text, icon.Text, group.Text, exec.Text)
	}, u.window)

	form.Resize(fyne.NewSize(dialogWidth, dialogHeight))
	form.Show()
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

	icon := newScrollFreeEntry()
	icon.SetText(command.Icon)
	icon.SetPlaceHolder("🚀")

	group := newScrollFreeEntry()
	group.SetText(command.Group)
	group.SetPlaceHolder("deploy")

	exec := newExecEntry()
	exec.SetText(command.Exec)

	description := widget.NewMultiLineEntry()
	description.SetText(command.Description)
	description.SetMinRowsVisible(2)
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
		widget.NewFormItem("Icon", u.iconField(icon)),
		widget.NewFormItem("Group", group),
		widget.NewFormItem("Exec", exec),
		widget.NewFormItem("Description", description),
		widget.NewFormItem("Params", params),
		widget.NewFormItem("Cwd", cwd),
		widget.NewFormItem("Shell", shell),
		widget.NewFormItem("Timeout", timeout),
	}

	form := dialog.NewForm("Edit command", "Save", "Cancel", items, func(ok bool) {
		if !ok {
			return
		}

		u.applyEdit(key, commandForm{
			Name:        name.Text,
			Aliases:     aliases.Text,
			Icon:        icon.Text,
			Group:       group.Text,
			Exec:        exec.Text,
			Description: description.Text,
			Params:      params.Text,
			Cwd:         cwd.Text,
			Shell:       shell.Text,
			Timeout:     timeout.Text,
		})
	}, u.window)

	form.Resize(fyne.NewSize(dialogWidth, editDialogHeight))
	form.Show()
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
