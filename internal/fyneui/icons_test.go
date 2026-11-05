//go:build fyne

package fyneui

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// The grid is the curated list plus the clear choice, and every emoji is
// unique.
func TestPickerEmojiList(t *testing.T) {
	if len(pickerEmojis) < 100 {
		t.Errorf("picker has %d emojis, want at least a hundred", len(pickerEmojis))
	}

	seen := make(map[string]bool)

	for _, emoji := range pickerEmojis {
		if emoji == "" {
			t.Error("an empty entry is in the grid")
		}

		if seen[emoji] {
			t.Errorf("%q appears twice", emoji)
		}

		seen[emoji] = true
	}
}

// Tapping an emoji fills the field and reports the choice; None clears it.
func TestIconGridSetsTheField(t *testing.T) {
	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})

	entry := newScrollFreeEntry()
	chosen := 0

	grid := front.iconGrid(entry, func() { chosen++ })

	if len(grid.Objects) != len(pickerEmojis)+1 {
		t.Fatalf("grid has %d cells, want %d", len(grid.Objects), len(pickerEmojis)+1)
	}

	first, ok := grid.Objects[1].(*widget.Button)
	if !ok {
		t.Fatalf("cell 1 = %T, want a button", grid.Objects[1])
	}

	test.Tap(first)

	if entry.Text != pickerEmojis[0] {
		t.Errorf("entry = %q, want %q", entry.Text, pickerEmojis[0])
	}

	if chosen != 1 {
		t.Errorf("chosen = %d, want the picker to close once", chosen)
	}

	none, ok := grid.Objects[0].(*widget.Button)
	if !ok {
		t.Fatalf("cell 0 = %T, want the clear button", grid.Objects[0])
	}

	test.Tap(none)

	if entry.Text != "" {
		t.Errorf("entry = %q, want it cleared", entry.Text)
	}
}
