//go:build fyne

package fyneui

import (
	"fyne.io/fyne/v2"

	"github.com/fiwon123/crower/internal/desktop"
)

// WindowStore remembers the window's size between runs.
//
// It is declared here, at the point of use: the window owns the conversion and
// the fallback, the command layer owns the file.
type WindowStore interface {
	// Load returns what was recorded, or ok=false when nothing usable was.
	Load() (desktop.WindowSize, bool)

	// Save records the size the window closes with.
	Save(size desktop.WindowSize) error
}

// initialSize is the size the window opens with: the remembered one when it is
// usable, the default otherwise.
func initialSize(store WindowStore) fyne.Size {
	size := desktop.DefaultWindowSize()

	if store != nil {
		if saved, ok := store.Load(); ok && saved.Usable() {
			size = saved
		}
	}

	return fyne.NewSize(float32(size.Width), float32(size.Height))
}

// recordedSize converts a canvas size, and says whether it is worth remembering.
func recordedSize(size fyne.Size) (desktop.WindowSize, bool) {
	recorded := desktop.WindowSize{Width: int(size.Width), Height: int(size.Height)}

	return recorded, recorded.Usable()
}
