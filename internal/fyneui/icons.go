//go:build fyne

package fyneui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// iconPickerSize is one cell of the emoji grid. A button with a single emoji
// needs enough room to render it at the theme's text size.
var iconPickerSize = fyne.NewSize(48, 48)

// pickerEmojis is the curated grid: a hundred markers a command might want,
// grouped so a reader can find one by theme rather than by scrolling a
// codepoint table. Adding one is a one-line change.
var pickerEmojis = buildPickerEmojis()

func buildPickerEmojis() []string {
	groups := [][]string{
		{
			// files and folders
			"📁", "📂", "🗂️", "📄", "📑", "🗄️", "📊", "📈", "📉", "📋",
			"📌", "🔖", "🧾", "📜", "📰", "📓", "📕", "📗", "📘", "📚",
		},
		{
			// writing and links
			"🔗", "📎", "📝", "✏️", "✒️", "🖊️", "🖋️", "🖌️", "✂️", "📐",
		},
		{
			// actions
			"▶️", "⏩", "⏭️", "⏮️", "⏪", "◀️", "🔁", "🔄", "➡️", "⬅️",
			"⬆️", "⬇️", "↗️", "↘️",
		},
		{
			// tools
			"🛠️", "🔧", "🔨", "⛏️", "🔩", "⚙️", "🧰", "🧲", "🧪", "🔬",
			"🔭", "🖥️", "💻", "⌨️", "🖱️", "🖨️", "📱", "📺", "📷", "🎥",
			"🔋", "🔌", "💡", "🔦",
		},
		{
			// status and events
			"✅", "❌", "⚠️", "ℹ️", "⏳", "⌛", "🚀", "🎯", "🔥", "💥",
			"✨", "⭐", "🌟", "🏁", "🎉", "🏆", "🥇",
		},
		{
			// weather and places
			"🌍", "🌎", "🌏", "🌕", "🌙", "☀️", "🌤️", "⛅", "☁️", "🌧️",
			"⛈️", "❄️", "☃️", "⛄", "💨", "🌈", "☔", "🌊",
		},
		{
			// animals
			"🐶", "🐱", "🐰", "🦊", "🐻", "🐼", "🐨", "🐯", "🦁", "🐮",
			"🐷", "🐸", "🐵", "🐧", "🦅", "🦉", "🐺", "🦄", "🐝", "🦋",
		},
	}

	seen := make(map[string]bool)
	emojis := make([]string, 0)

	for _, group := range groups {
		for _, emoji := range group {
			if seen[emoji] {
				continue
			}

			seen[emoji] = true
			emojis = append(emojis, emoji)
		}
	}

	return emojis
}

// showIconPicker opens the emoji grid over the window. Choosing one fills the
// field and closes the picker; the field stays editable for anything not in the
// grid.
func (u *ui) showIconPicker(target *widget.Entry) {
	var picker *dialog.CustomDialog

	grid := u.iconGrid(target, func() {
		if picker != nil {
			picker.Dismiss()
		}
	})

	picker = dialog.NewCustom("Pick an icon", "Close", container.NewVScroll(grid), u.window)
	picker.Resize(fyne.NewSize(dialogWidth, dialogHeight))
	picker.Show()
}

// iconGrid is the picker's content without the dialog, so a test can tap it.
// The first button clears the field; every other one sets it.
func (u *ui) iconGrid(target *widget.Entry, chosen func()) *fyne.Container {
	grid := container.NewGridWrap(iconPickerSize)

	grid.Add(widget.NewButton("None", func() {
		target.SetText("")
		chosen()
	}))

	for _, emoji := range pickerEmojis {
		emoji := emoji

		grid.Add(widget.NewButton(emoji, func() {
			target.SetText(emoji)
			chosen()
		}))
	}

	return grid
}

// iconField is the Icon row: the entry with the picker beside it.
func (u *ui) iconField(target *widget.Entry) fyne.CanvasObject {
	pick := widget.NewButtonWithIcon("Pick…", theme.GridIcon(), func() {
		u.showIconPicker(target)
	})
	pick.Importance = widget.LowImportance

	return container.NewBorder(nil, nil, nil, pick, target)
}
