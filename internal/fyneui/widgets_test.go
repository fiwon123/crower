//go:build fyne

package fyneui

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

// growingRows counts explicit lines and wrapped rows, between the bounds the
// field was built with.
func TestGrowingRows(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		minRows int
		maxRows int
		want    int
	}{
		{name: "empty", text: "", minRows: 4, maxRows: 12, want: 4},
		{name: "one short line", text: "echo hi", minRows: 4, maxRows: 12, want: 4},
		{name: "six lines", text: "a\nb\nc\nd\ne\nf", minRows: 4, maxRows: 12, want: 6},
		{name: "capped", text: strings.Repeat("line\n", 20), minRows: 4, maxRows: 12, want: 12},
		{name: "one wrapped line", text: strings.Repeat("x", 150), minRows: 4, maxRows: 12, want: 4},
		{name: "one long wrapped line", text: strings.Repeat("x", 300), minRows: 4, maxRows: 12, want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := growingRows(tt.text, tt.minRows, tt.maxRows); got != tt.want {
				t.Errorf("growingRows(%q) = %d, want %d", tt.text, got, tt.want)
			}
		})
	}
}

// Setting text on a growing field raises its minimum height, which is what the
// dialog lays out; the cap keeps it from taking the window.
func TestGrowingEntryGrowsWithItsText(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	entry := newGrowingEntry(4, 10)
	entry.Resize(fyne.NewSize(600, 120))

	short := entry.MinSize().Height

	entry.SetText("1\n2\n3\n4\n5\n6\n7\n8")

	tall := entry.MinSize().Height

	if tall <= short {
		t.Errorf("MinSize height = %v, want it to grow past %v", tall, short)
	}

	entry.SetText(strings.Repeat("line\n", 40))

	capped := entry.MinSize().Height

	if got := growingRows(entry.Text, 4, 10); got != 10 {
		t.Fatalf("growingRows = %d, want the cap", got)
	}

	if capped < tall {
		t.Errorf("capped MinSize height = %v, want at least %v", capped, tall)
	}
}
