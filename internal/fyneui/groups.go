//go:build fyne

package fyneui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

// showGroups edits the icon of every group in use.
func (u *ui) showGroups() {
	groups := u.model.Groups()
	if len(groups) == 0 {
		showText(
			"Groups",
			"No command belongs to a group yet. Set one in the create or edit dialog, then come back.",
			u.window,
		)

		return
	}

	entries := make(map[string]*iconChoice, len(groups))
	items := make([]*widget.FormItem, 0, len(groups))

	for _, name := range groups {
		choice, field := u.iconField(u.groupIcons[name])
		entries[name] = choice

		items = append(items, widget.NewFormItem(name, field))
	}

	form := container.NewVScroll(widget.NewForm(items...))
	form.SetMinSize(fyne.NewSize(420, 320))

	dialog.ShowCustomConfirm("Group icons", "Save", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		icons := make(map[string]string, len(entries))
		for name, choice := range entries {
			icons[name] = choice.value
		}

		u.applyGroupIcons(icons)
	}, u.window)
}

// applyGroupIcons stores what the dialog holds and refreshes the filter. It is
// synchronous, so a test can drive it without the popup.
func (u *ui) applyGroupIcons(icons map[string]string) bool {
	if u.groups == nil {
		return false
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
