//go:build fyne

package fyneui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// groupLabel is what a group is called in the window: its icon, when it has
// one, then its name.
func groupLabel(icon, name string) string {
	if icon == "" {
		return name
	}

	return icon + " " + name
}

// showGroups edits the groups the window offers: one compact row each, with the
// icon button, the name and a Delete action, plus New group.
//
// A group is stored on its own, so it can be created here before any command
// uses it. The dialog is sized from the window rather than given a fixed size:
// a row list in a fixed 420x320 box wasted the space of a maximised window and
// could not fit a tiny one.
func (u *ui) showGroups() {
	entries := make(map[string]*iconChoice)
	list := container.NewVBox()

	empty := widget.NewLabel("No groups yet. Create one here, then pick it in a command.")
	empty.Wrapping = fyne.TextWrapWord
	empty.Selectable = true

	updateEmpty := func() {
		if len(entries) == 0 {
			empty.Show()
			return
		}

		empty.Hide()
	}

	var addRow func(name string)

	addRow = func(name string) {
		if _, exists := entries[name]; exists {
			return
		}

		choice, field := u.iconField(u.groupIcons[name])
		entries[name] = choice

		nameLabel := widget.NewLabel(name)
		nameLabel.Selectable = true

		var row fyne.CanvasObject

		remove := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			confirmNamed("Delete group", fmt.Sprintf(
				"Delete the group %q?\n\nA command that names it keeps the name until you change the command.",
				name,
			), name, "Delete", widget.DangerImportance, u.window, func(confirmed bool) {
				if !confirmed || !u.removeGroup(name) {
					return
				}

				delete(entries, name)
				list.Remove(row)
				updateEmpty()
				list.Refresh()
			})
		})
		remove.Importance = widget.DangerImportance

		row = container.NewHBox(field, nameLabel, layout.NewSpacer(), remove, gap(theme.Padding()*2, 0))
		list.Add(row)

		updateEmpty()
	}

	for _, name := range u.groupNames() {
		addRow(name)
	}

	updateEmpty()

	newGroup := widget.NewButtonWithIcon("New group", theme.ContentAddIcon(), func() {
		u.promptNewGroup(func(name, icon string) {
			if !u.createGroup(name, icon) {
				return
			}

			addRow(name)
			entries[name].button.SetEmoji(strings.TrimSpace(icon))
			list.Refresh()
		})
	})

	newGroup.Importance = widget.HighImportance

	content := container.NewBorder(
		container.NewVBox(newGroup, gap(0, theme.Padding()), widget.NewSeparator()),
		nil,
		nil,
		nil,
		container.NewVScroll(list),
	)

	showForm("Groups", "Save", "Cancel", content, u.window, func(ok bool) {
		if !ok {
			return
		}

		icons := make(map[string]string, len(entries))
		for name, choice := range entries {
			icons[name] = choice.value
		}

		u.applyGroupIcons(icons)
	}, dialogSizeFor(u.window, 380, 260, 720, 620))
}

// createGroup stores a group and its icon, so it is offered in the filter and in
// a command's dropdown without one using it yet. It is synchronous, so a test
// can drive it without the popup.
func (u *ui) createGroup(name, icon string) bool {
	name = strings.TrimSpace(name)
	icon = strings.TrimSpace(icon)

	if name == "" {
		showError("Enter a group name.", u.window)
		return false
	}

	if u.groups == nil {
		u.groupIcons[name] = icon
		u.refreshGroups()

		return true
	}

	icons, err := u.groups.SetIcon(name, icon)
	if err != nil {
		showError(err.Error(), u.window)
		return false
	}

	u.groupIcons = icons
	u.refreshGroups()

	return true
}

// removeGroup deletes a stored group. A group a command still names comes back
// without an icon, because the catalog derives it from the commands.
func (u *ui) removeGroup(name string) bool {
	if u.groups == nil {
		delete(u.groupIcons, name)
		u.refreshGroups()

		return true
	}

	icons, err := u.groups.Remove(name)
	if err != nil {
		showError(err.Error(), u.window)
		return false
	}

	u.groupIcons = icons
	u.refreshGroups()

	return true
}

// applyGroupIcons stores what the dialog holds and refreshes the filter. It is
// synchronous, so a test can drive it without the popup.
func (u *ui) applyGroupIcons(icons map[string]string) bool {
	if u.groups == nil {
		for name, icon := range icons {
			u.groupIcons[name] = strings.TrimSpace(icon)
		}

		u.refreshGroups()

		return true
	}

	for name, icon := range icons {
		updated, err := u.groups.SetIcon(name, strings.TrimSpace(icon))
		if err != nil {
			showError(err.Error(), u.window)
			return false
		}

		u.groupIcons = updated
	}

	u.refreshGroups()

	return true
}
