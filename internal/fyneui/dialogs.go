//go:build fyne

package fyneui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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

	group, groupWidget := u.groupField("")

	exec := newExecEntry()

	// The form scrolls, so a grown exec field or a short window never clips the
	// last row or the buttons.
	form := newFormGrid([]formRow{
		{"Name", name},
		{"Aliases", aliases},
		{"Icon", iconWidget},
		{"Group", groupWidget},
		{"Exec", exec},
	})

	showForm("New command", "Create", "Cancel", container.NewVScroll(form), u.window, func(ok bool) {
		if !ok {
			return
		}

		u.applyCreate(name.Text, aliases.Text, icon.value, group.value, exec.Text)
	}, fyne.NewSize(dialogWidth, dialogHeight))
}

// noGroupLabel is the group select's empty choice.
const noGroupLabel = "No group"

// groupChoice is the group field's value and selector. The selector's label
// carries the group's icon, so the two maps translate between the stored name
// and what the row shows.
type groupChoice struct {
	value        string
	labels       map[string]string
	names        map[string]string
	selectWidget *widget.Select
}

// groupField lets a command pick an existing group or create one, instead of
// typing a name. A created group is stored immediately, so it is offered here
// and in the filter before any command uses it.
func (u *ui) groupField(initial string) (*groupChoice, fyne.CanvasObject) {
	choice := &groupChoice{
		value:  strings.TrimSpace(initial),
		labels: make(map[string]string),
		names:  make(map[string]string),
	}

	options := []string{noGroupLabel}

	for _, name := range u.groupNames() {
		label := groupLabel(u.groupIcons[name], name)
		choice.labels[name] = label
		choice.names[label] = name
		options = append(options, label)
	}

	// A value set outside the window (`create --group`) must still show, even
	// when no command in the list carries it any more.
	if choice.value != "" && choice.labels[choice.value] == "" {
		label := groupLabel(u.groupIcons[choice.value], choice.value)
		choice.labels[choice.value] = label
		choice.names[label] = choice.value
		options = append(options, label)
	}

	choice.selectWidget = widget.NewSelect(options, func(label string) {
		if label == noGroupLabel {
			choice.value = ""
			return
		}

		choice.value = choice.names[label]
	})
	choice.selectWidget.SetSelected(choice.groupLabel())

	newGroup := newButton("New group", nil, widget.MediumImportance, func() {
		u.promptNewGroup(func(name, icon string) {
			if !u.createGroup(name, icon) {
				return
			}

			label := groupLabel(icon, name)
			choice.labels[name] = label
			choice.names[label] = name
			choice.selectWidget.Options = append(choice.selectWidget.Options, label)
			choice.selectWidget.Refresh()
			choice.selectWidget.SetSelected(label)
			choice.value = name
		})
	})

	return choice, container.NewBorder(nil, nil, nil, newGroup, choice.selectWidget)
}

func (c *groupChoice) groupLabel() string {
	if c.value == "" {
		return noGroupLabel
	}

	if label, found := c.labels[c.value]; found {
		return label
	}

	return c.value
}

// promptNewGroup asks for a group's name and icon and hands both to onCreate.
// A group is a name first; the icon is what the filter shows beside it.
func (u *ui) promptNewGroup(onCreate func(name, icon string)) {
	name := newScrollFreeEntry()
	name.SetPlaceHolder("deploy")

	icon, iconWidget := u.iconField("")

	form := newFormGrid([]formRow{
		{"Name", name},
		{"Icon", iconWidget},
	})

	showForm("New group", "Create", "Cancel", form, u.window, func(ok bool) {
		if !ok {
			return
		}

		value := strings.TrimSpace(name.Text)
		if value == "" {
			showError("Enter a group name.", u.window)
			return
		}

		onCreate(value, icon.value)
	}, fyne.NewSize(460, 260))
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

	group, groupWidget := u.groupField(command.Group)

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

	form := newFormGrid([]formRow{
		{"Name", name},
		{"Aliases", aliases},
		{"Icon", iconWidget},
		{"Group", groupWidget},
		{"Exec", exec},
		{"Description", description},
		{"Params", params},
		{"Cwd", cwd},
		{"Shell", shell},
		{"Timeout", timeout},
	})

	showForm("Edit command", "Save", "Cancel", container.NewVScroll(form), u.window, func(ok bool) {
		if !ok {
			return
		}

		u.applyEdit(key, commandForm{
			Name:        name.Text,
			Aliases:     aliases.Text,
			Icon:        icon.value,
			Group:       group.value,
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

	confirmNamed("Delete command", question, key, "Delete", widget.DangerImportance, u.window, func(ok bool) {
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

	// showForm, not Fyne's native dialog: the native one puts Run and Cancel at
	// the right edge, so the prompt looked different from every other dialog
	// (#500).
	showForm("Parameters", "Run", "Cancel", container.NewVBox(input), u.window, func(ok bool) {
		if !ok {
			u.model.State = desktop.ViewList
			return
		}

		u.model.ParamInput = input.Text
		u.model.SubmitParams()
		u.showOutput()
	}, dialogSizeFor(u.window, 420, 200, 620, 300))
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
	closeButton := newButton("Close", nil, widget.MediumImportance, nil)

	body := dialogScaffold(title, content, dialogFooter(closeButton))

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

	confirmButton := newButton(confirm, nil, widget.HighImportance, func() {
		pop.Hide()
		onConfirm(true)
	})

	dismissButton := newButton(dismiss, nil, widget.LowImportance, func() {
		if pop != nil {
			pop.Hide()
		}

		onConfirm(false)
	})

	body := dialogScaffold(title, content, dialogFooter(dismissButton, confirmButton))

	pop = widget.NewPopUp(container.NewPadded(body), window.Canvas())
	pop.Resize(size)
	pop.Show()

	return pop
}

// dialogSizeFor fits a dialog to the window: a margin at each side, but never
// smaller than min and never larger than max. A tiny window gets a dialog that
// still fits on it; a maximised one does not get a dialog stretched across the
// screen with a single column of text in the middle.
func dialogSizeFor(window fyne.Window, minWidth, minHeight, maxWidth, maxHeight float32) fyne.Size {
	size := window.Canvas().Size()
	pad := theme.Padding() * 6

	width := max(minWidth, min(maxWidth, size.Width-pad))
	height := max(minHeight, min(maxHeight, size.Height-pad))

	return fyne.NewSize(width, height)
}

// dialogFooter is a dialog's button row: centred, with air under it so the
// buttons are not glued to the window's lower edge.
func dialogFooter(buttons ...fyne.CanvasObject) fyne.CanvasObject {
	return container.NewVBox(
		container.NewCenter(container.NewHBox(buttons...)),
		gap(0, theme.Padding()*3),
	)
}

// dialogScaffold is a dialog's body: a title over a separator, the content, and
// the button row beneath it. Every dialog and the screenshot preview build it
// with this function, so a capture cannot show a footer, or a button order, that
// users never get (#503).
func dialogScaffold(title string, content, footer fyne.CanvasObject) fyne.CanvasObject {
	head := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	return container.NewBorder(
		container.NewVBox(head, widget.NewSeparator()),
		footer,
		nil,
		nil,
		content,
	)
}

// confirmNamed asks for a destructive action to be confirmed by typing the
// name. A click on a confirm button is one slip away; typing the name is not, so
// the action stays disabled until the field matches exactly.
func confirmNamed(title, message, expected, action string, importance widget.Importance, window fyne.Window, onConfirm func(bool)) *widget.PopUp {
	entry := newScrollFreeEntry()
	entry.SetPlaceHolder(expected)

	hint := newSelectableText(fmt.Sprintf("Type %s to %s.", expected, strings.ToLower(action)))
	messageLabel := newSelectableText(message)

	var pop *widget.PopUp

	confirm := newButton(action, nil, importance, func() {
		pop.Hide()
		onConfirm(true)
	})
	confirm.Disable()

	dismiss := newButton("Cancel", nil, widget.LowImportance, func() {
		if pop != nil {
			pop.Hide()
		}
	})

	entry.OnChanged = func(text string) {
		if strings.TrimSpace(text) == expected {
			confirm.Enable()
			return
		}

		confirm.Disable()
	}

	body := dialogScaffold(title, container.NewVBox(messageLabel, hint, entry), dialogFooter(dismiss, confirm))

	pop = widget.NewPopUp(container.NewPadded(body), window.Canvas())
	pop.Resize(dialogSizeFor(window, 420, 240, 620, 440))
	pop.Show()

	return pop
}

// showConfirm asks a question. It stays modal so an outside click cannot answer
// it, and its body is a selectable label so the message can be copied.
//
// It is built here rather than with dialog.NewCustomConfirm because Fyne's
// native dialog puts its buttons at the right edge; the rest of the window
// centres them (#500).
func showConfirm(title, message string, window fyne.Window, onConfirm func(bool)) *widget.PopUp {
	var pop *widget.PopUp

	confirm := newButton("Yes", nil, widget.HighImportance, func() {
		pop.Hide()
		onConfirm(true)
	})

	dismiss := newButton("No", nil, widget.LowImportance, func() {
		pop.Hide()
		onConfirm(false)
	})

	body := dialogScaffold(title, newSelectableText(message), dialogFooter(dismiss, confirm))

	pop = widget.NewModalPopUp(container.NewPadded(body), window.Canvas())
	pop.Resize(dialogSizeFor(window, 420, 220, 620, 380))
	pop.Show()

	return pop
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
