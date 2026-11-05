package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/pkg/utils"
)

// windowName is the file that remembers the window's size, beside the
// configuration.
const windowName = "window.toml"

// WindowState is the on-disk shape of the remembered window size.
type WindowState struct {
	Width  int `toml:"width"`
	Height int `toml:"height"`
}

// Window stores the size the window had when it last closed.
//
// Like groups.toml it is presentation the CLI never reads or writes: a missing
// file means "open at the default size", not an error.
type Window struct {
	path string
}

// NewWindow addresses the remembered size of the installation stored in folder.
func NewWindow(folder string) *Window {
	return &Window{path: filepath.Join(folder, windowName)}
}

// Path is where the remembered size lives.
func (w *Window) Path() string {
	return w.path
}

// Load reads the recorded size. A missing file is not an error; ok is false.
func (w *Window) Load() (WindowState, bool, error) {
	state := WindowState{}

	if _, err := os.Stat(w.path); errors.Is(err, os.ErrNotExist) {
		return state, false, nil
	}

	if err := utils.ReadToml(w.path, &state); err != nil {
		return WindowState{}, false, fmt.Errorf("read window size %s: %w", w.path, err)
	}

	return state, state.Width > 0 && state.Height > 0, nil
}

// Save writes the size.
func (w *Window) Save(state WindowState) error {
	if err := utils.WriteToml(state, w.path); err != nil {
		return fmt.Errorf("write window size %s: %w", w.path, err)
	}

	return nil
}
