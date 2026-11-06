//go:build fyne

package fyneui

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed emoji/emoji.json
var emojiJSON []byte

// emojiCategories is the picker's category order. The catalogue's own order is
// data-dependent, so the familiar groups are listed explicitly.
var emojiCategories = []string{
	"Smileys & Emotion",
	"People & Body",
	"Animals & Nature",
	"Food & Drink",
	"Travel & Places",
	"Activities",
	"Objects",
	"Symbols",
	"Flags",
}

// emojiEntry is one row of the gemoji catalogue.
type emojiEntry struct {
	Emoji       string   `json:"emoji"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Aliases     []string `json:"aliases"`
	Tags        []string `json:"tags"`
}

var (
	emojiOnce      sync.Once
	emojiCatalogue []emojiEntry
)

// loadEmojiCatalogue parses the embedded data once, when the picker first opens.
// An entry without a glyph is skipped so every cell is tappable.
func loadEmojiCatalogue() []emojiEntry {
	emojiOnce.Do(func() {
		var raw []emojiEntry
		if err := json.Unmarshal(emojiJSON, &raw); err != nil {
			fyne.LogError("loading the emoji catalogue", err)
			return
		}

		emojiCatalogue = make([]emojiEntry, 0, len(raw))

		for _, entry := range raw {
			if entry.Emoji == "" {
				continue
			}

			emojiCatalogue = append(emojiCatalogue, entry)
		}
	})

	return emojiCatalogue
}

// filterEmoji returns entries whose name, alias or tag contains query, and
// whose category matches (or "All"). The match is case-insensitive.
func filterEmoji(entries []emojiEntry, query, category string) []emojiEntry {
	query = strings.ToLower(strings.TrimSpace(query))
	filtered := make([]emojiEntry, 0, len(entries))

	for _, entry := range entries {
		if category != "" && category != "All" && entry.Category != category {
			continue
		}

		if query != "" && !emojiMatches(entry, query) {
			continue
		}

		filtered = append(filtered, entry)
	}

	return filtered
}

func emojiMatches(entry emojiEntry, query string) bool {
	query = strings.ToLower(strings.TrimSpace(query))

	if strings.Contains(strings.ToLower(entry.Description), query) {
		return true
	}

	if strings.Contains(entry.Emoji, query) {
		return true
	}

	for _, alias := range entry.Aliases {
		if strings.Contains(strings.ToLower(alias), query) {
			return true
		}
	}

	for _, tag := range entry.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}

	return false
}

// emojiColumns is how many emoji a picker row holds.
const emojiColumns = 8

// emojiTextSize is the glyph size in the picker. The theme's default text size
// makes an emoji a tiny dot in a grid cell.
const emojiTextSize = 30

var emojiCellSize = fyne.NewSize(56, 56)

// emojiCell is one picker cell: a large glyph centred in a fixed box. The table
// handles the tap, so the cell itself does not need to be a button.
func emojiCell() fyne.CanvasObject {
	text := canvas.NewText("", theme.Color(theme.ColorNameForeground))
	text.TextSize = emojiTextSize
	text.Alignment = fyne.TextAlignCenter

	return container.NewGridWrap(emojiCellSize, text)
}

// emojiRows is the table row count for n entries.
func emojiRows(n int) int {
	if n == 0 {
		return 0
	}

	return (n + emojiColumns - 1) / emojiColumns
}

// showEmojiPicker opens the catalogue over the window. Choosing an emoji calls
// onPick and closes the pop-up; clicking outside closes it too.
//
// The grid is a widget.Table, so only the visible cells exist: the catalogue is
// nearly two thousand entries, and a plain grid would create a cell for every
// one of them.
func (u *ui) showEmojiPicker(current string, onPick func(string)) {
	entries := loadEmojiCatalogue()
	filtered := filterEmoji(entries, "", "All")

	search := newScrollFreeEntry()
	search.SetPlaceHolder("Search by name, tag or emoji")

	category := widget.NewSelect(append([]string{"All"}, emojiCategories...), nil)
	category.SetSelected("All")

	count := widget.NewLabel("")
	count.Alignment = fyne.TextAlignCenter

	var pop *widget.PopUp

	updateCell := func(id widget.TableCellID, object fyne.CanvasObject) {
		cell, ok := object.(*fyne.Container)
		if !ok || len(cell.Objects) == 0 {
			return
		}

		text, ok := cell.Objects[0].(*canvas.Text)
		if !ok {
			return
		}

		index := id.Row*emojiColumns + id.Col
		if index >= len(filtered) {
			text.Text = ""
			text.Refresh()

			return
		}

		text.Text = filtered[index].Emoji
		text.Refresh()
	}

	table := widget.NewTable(
		func() (int, int) { return emojiRows(len(filtered)), emojiColumns },
		emojiCell,
		updateCell,
	)
	table.HideSeparators = true
	table.OnSelected = func(id widget.TableCellID) {
		index := id.Row*emojiColumns + id.Col
		if index >= len(filtered) {
			return
		}

		onPick(filtered[index].Emoji)

		if pop != nil {
			pop.Hide()
		}
	}

	for column := 0; column < emojiColumns; column++ {
		table.SetColumnWidth(column, 60)
	}

	refilter := func() {
		filtered = filterEmoji(entries, search.Text, category.Selected)
		count.SetText(fmt.Sprintf("%d emoji", len(filtered)))
		table.Refresh()
		table.ScrollToTop()
	}

	search.OnChanged = func(string) { refilter() }
	category.OnChanged = func(string) { refilter() }
	refilter()

	top := container.NewVBox(
		search,
		container.NewBorder(nil, nil, widget.NewLabel("Category"), nil, category),
		count,
	)
	content := container.NewBorder(top, nil, nil, nil, table)

	title := "Pick an icon"
	if current != "" {
		title += " — " + current
	}

	pop = showDismissible(title, content, u.window, fyne.NewSize(560, 620))
}

// iconChoice is the read-only icon value a form holds between the picker and
// its submit.
type iconChoice struct {
	value string
	label *widget.Label
}

// iconDisplay is what the field shows: the emoji, or "None".
func iconDisplay(value string) string {
	if strings.TrimSpace(value) == "" {
		return "None"
	}

	return value
}

// iconField builds the Icon row: a read-only display with Pick… and Clear, so
// the value can only come from the catalogue.
func (u *ui) iconField(initial string) (*iconChoice, fyne.CanvasObject) {
	choice := &iconChoice{value: strings.TrimSpace(initial)}
	choice.label = widget.NewLabel(iconDisplay(choice.value))

	pick := widget.NewButtonWithIcon("Pick…", theme.GridIcon(), func() {
		u.showEmojiPicker(choice.value, func(emoji string) {
			choice.value = emoji
			choice.label.SetText(iconDisplay(emoji))
		})
	})
	pick.Importance = widget.LowImportance

	clear := widget.NewButton("Clear", func() {
		choice.value = ""
		choice.label.SetText(iconDisplay(""))
	})

	return choice, container.NewBorder(nil, nil, nil, container.NewHBox(clear, pick), choice.label)
}
