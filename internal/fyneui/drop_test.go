//go:build fyne

package fyneui

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// A drop appends to the Sources field, opens the Tools tab and says how many
// paths arrived — without running anything.
func TestDropFillsTheTransferSources(t *testing.T) {
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	window.SetContent(front.content())
	window.Resize(fyneSize())

	front.sourcesEntry.SetText("/tmp/keep.txt")

	first, err := storage.ParseURI("file:///tmp/a.txt")
	if err != nil {
		t.Fatalf("ParseURI: %v", err)
	}

	second, err := storage.ParseURI("file:///tmp/b.txt")
	if err != nil {
		t.Fatalf("ParseURI: %v", err)
	}

	front.handleDrop(fyne.NewPos(0, 0), []fyne.URI{first, second})

	want := "/tmp/keep.txt, /tmp/a.txt, /tmp/b.txt"
	if front.sourcesEntry.Text != want {
		t.Errorf("sources = %q, want %q", front.sourcesEntry.Text, want)
	}

	if front.tabs.Selected() == nil || front.tabs.Selected().Text != "Tools" {
		t.Errorf("selected tab = %v, want Tools", front.tabs.Selected())
	}

	if !strings.Contains(front.transferResult.label.Text, "Added 2 path(s)") {
		t.Errorf("result = %q, want the drop count", front.transferResult.label.Text)
	}
}

// A drop with nothing usable changes nothing and stays where it was.
func TestEmptyDropDoesNothing(t *testing.T) {
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	window.SetContent(front.content())
	window.Resize(fyneSize())

	front.sourcesEntry.SetText("/tmp/keep.txt")
	front.tabs.SelectIndex(0)
	front.handleDrop(fyne.NewPos(0, 0), nil)

	if front.sourcesEntry.Text != "/tmp/keep.txt" {
		t.Errorf("sources = %q, want them untouched", front.sourcesEntry.Text)
	}

	if front.tabs.Selected().Text != "Commands" {
		t.Errorf("selected tab = %q, want Commands kept", front.tabs.Selected().Text)
	}
}
