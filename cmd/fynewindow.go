package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/desktop"
)

// fyneWindow remembers the window's size in window.toml beside the
// configuration.
//
// It has no Fyne import, so the adapter is covered by the normal cmd tests.
type fyneWindow struct {
	folder string
}

func newFyneWindow(folder string) fyneWindow {
	return fyneWindow{folder: folder}
}

// Load returns the recorded size when it is usable.
func (w fyneWindow) Load() (desktop.WindowSize, bool) {
	state, ok, err := storage.NewWindow(w.folder).Load()
	if err != nil || !ok {
		return desktop.WindowSize{}, false
	}

	size := desktop.WindowSize{Width: state.Width, Height: state.Height}

	return size, size.Usable()
}

// Save records the size, ignoring one that could not open a usable window: a
// window dragged to 200x100 should not come back that way.
func (w fyneWindow) Save(size desktop.WindowSize) error {
	if !size.Usable() {
		return nil
	}

	return storage.NewWindow(w.folder).Save(storage.WindowState{Width: size.Width, Height: size.Height})
}
