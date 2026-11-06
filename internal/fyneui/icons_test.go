//go:build fyne

package fyneui

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// The embedded catalogue is the full gemoji set, with every entry usable.
func TestEmojiCatalogue(t *testing.T) {
	entries := loadEmojiCatalogue()

	if len(entries) < 1800 {
		t.Fatalf("catalogue = %d entries, want the full gemoji set", len(entries))
	}

	categories := make(map[string]bool)
	seen := make(map[string]bool)

	for _, entry := range entries {
		if entry.Emoji == "" {
			t.Fatal("the catalogue has an entry without a glyph")
		}

		if seen[entry.Emoji] {
			t.Errorf("%q appears twice", entry.Emoji)
		}

		seen[entry.Emoji] = true
		categories[entry.Category] = true
	}

	for _, want := range emojiCategories {
		if !categories[want] {
			t.Errorf("category %q is missing from the catalogue", want)
		}
	}
}

// The picker filters by name, alias or tag, and by category.
func TestFilterEmoji(t *testing.T) {
	entries := loadEmojiCatalogue()

	all := filterEmoji(entries, "", "All")
	if len(all) != len(entries) {
		t.Errorf("all = %d, want %d", len(all), len(entries))
	}

	smiles := filterEmoji(entries, "SMILE", "All")
	if len(smiles) == 0 {
		t.Error("query SMILE matched nothing")
	}

	found := false
	for _, entry := range smiles {
		if entry.Emoji == "😀" {
			found = true
		}
	}

	if !found {
		t.Error("query SMILE did not match the grinning face")
	}

	rocket := filterEmoji(entries, "🚀", "All")
	if len(rocket) == 0 {
		t.Error("searching for the emoji glyph matched nothing")
	}

	flagMatches := filterEmoji(entries, "flag", "Flags")
	if len(flagMatches) == 0 {
		t.Error("category+query matched nothing")
	}

	for _, entry := range flagMatches {
		if entry.Category != "Flags" || !emojiMatches(entry, "flag") {
			t.Errorf("%q does not satisfy the category and query", entry.Emoji)
		}
	}

	flags := filterEmoji(entries, "", "Flags")
	if len(flags) == 0 {
		t.Fatal("the Flags category is empty")
	}

	for _, entry := range flags {
		if entry.Category != "Flags" {
			t.Errorf("%q is in %q, want Flags", entry.Emoji, entry.Category)
		}
	}

	if got := filterEmoji(entries, "definitely-not-an-emoji", "All"); len(got) != 0 {
		t.Errorf("nonsense query = %d entries, want none", len(got))
	}

	if got := emojiRows(0); got != 0 {
		t.Errorf("emojiRows(0) = %d, want 0", got)
	}

	if got := emojiRows(emojiColumns + 1); got != 2 {
		t.Errorf("emojiRows(%d) = %d, want 2", emojiColumns+1, got)
	}
}

// The icon field is a read-only display with Clear and Pick…; it is not an
// entry, so an arbitrary string cannot be typed in.
func TestIconFieldIsReadOnly(t *testing.T) {
	front, window := newTestUI(t, testModel(), stubEditor{}, stubReverter{})

	choice, field := front.iconField("🚀")

	if choice.label.Text != "🚀" {
		t.Errorf("display = %q, want the emoji", choice.label.Text)
	}

	if hasEntry(field) {
		t.Error("the icon field contains an editable entry")
	}

	clear := findButton(field, "Clear")
	if clear == nil {
		t.Fatal("the field has no Clear button")
	}

	clear.OnTapped()

	if choice.value != "" || choice.label.Text != "None" {
		t.Errorf("cleared field = %q/%q, want None", choice.value, choice.label.Text)
	}

	pick := findButton(field, "Pick…")
	if pick == nil {
		t.Fatal("the field has no Pick… button")
	}

	pick.OnTapped()

	if len(window.Canvas().Overlays().List()) == 0 {
		t.Error("Pick… did not open the catalogue")
	}
}

// hasEntry reports whether an editable entry is reachable.
func hasEntry(object fyne.CanvasObject) bool {
	switch typed := object.(type) {
	case *widget.Entry:
		return true
	case *fyne.Container:
		for _, child := range typed.Objects {
			if hasEntry(child) {
				return true
			}
		}
	case fyne.Widget:
		renderer := test.WidgetRenderer(typed)
		if renderer == nil {
			return false
		}

		for _, child := range renderer.Objects() {
			if hasEntry(child) {
				return true
			}
		}
	}

	return false
}

// findButton returns the button with exactly this text, or nil.
func findButton(object fyne.CanvasObject, text string) *widget.Button {
	switch typed := object.(type) {
	case *widget.Button:
		if strings.TrimSpace(typed.Text) == text {
			return typed
		}
	case *fyne.Container:
		for _, child := range typed.Objects {
			if button := findButton(child, text); button != nil {
				return button
			}
		}
	case fyne.Widget:
		renderer := test.WidgetRenderer(typed)
		if renderer == nil {
			return nil
		}

		for _, child := range renderer.Objects() {
			if button := findButton(child, text); button != nil {
				return button
			}
		}
	}

	return nil
}
