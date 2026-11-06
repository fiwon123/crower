//go:build fyne

package fyneui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// showNewSetup asks for a name and creates an empty setup, the way
// `crower setup new` does.
func (u *ui) showNewSetup() {
	if u.setups == nil {
		return
	}

	name := newScrollFreeEntry()
	name.SetPlaceHolder("streamer")

	form := widget.NewForm(widget.NewFormItem("Name", name))

	showForm("New setup", "Create", "Cancel", form, u.window, func(ok bool) {
		if !ok {
			return
		}

		u.createSetup(name.Text)
	}, fyne.NewSize(dialogWidth, dialogHeight))
}

// createSetup is the form's submit path.
func (u *ui) createSetup(name string) bool {
	if u.setups == nil {
		return false
	}

	name = strings.TrimSpace(name)
	if name == "" {
		showError("Enter a setup name.", u.window)
		return false
	}

	if err := u.setups.Create(name); err != nil {
		showError(err.Error(), u.window)
		return false
	}

	u.refreshSetupSelect()

	return true
}

// showDeleteSetup asks before removing the setup the selector points at.
//
// The port refuses the active setup, which is the one the window has loaded, so
// the confirmation says so rather than offering a choice that would fail.
func (u *ui) showDeleteSetup() {
	if u.setups == nil || u.setupSelect == nil {
		return
	}

	name := u.setupSelect.Selected
	if name == "" {
		return
	}

	question := fmt.Sprintf(
		"Delete setup %q and its history?\n\nIts commands and its history are removed. The default setup cannot be deleted, and neither can the setup the window is using.",
		name,
	)

	showConfirm("Delete setup", question, u.window, func(ok bool) {
		if !ok {
			return
		}

		u.deleteSetup(name)
	})
}

// deleteSetup is the confirmation's submit path.
func (u *ui) deleteSetup(name string) bool {
	if u.setups == nil {
		return false
	}

	if err := u.setups.Delete(name); err != nil {
		showError(err.Error(), u.window)
		return false
	}

	u.refreshSetupSelect()

	return true
}
