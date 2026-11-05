package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage"
)

// A missing file is "no size recorded", not an error.
func TestWindowLoadIsEmptyBeforeAnythingIsSaved(t *testing.T) {
	window := storage.NewWindow(t.TempDir())

	_, ok, err := window.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if ok {
		t.Error("a fresh folder reported a recorded size")
	}
}

func TestWindowSaveLoadRoundTrip(t *testing.T) {
	folder := t.TempDir()
	window := storage.NewWindow(folder)

	if err := window.Save(storage.WindowState{Width: 1200, Height: 800}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	state, ok, err := window.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if !ok {
		t.Fatal("the saved size was not reported")
	}

	if state.Width != 1200 || state.Height != 800 {
		t.Errorf("state = %+v, want 1200x800", state)
	}

	if _, err := os.Stat(filepath.Join(folder, "window.toml")); err != nil {
		t.Errorf("window.toml was not written: %v", err)
	}
}

// A file with no usable values reports "nothing recorded".
func TestWindowLoadRejectsEmptyValues(t *testing.T) {
	folder := t.TempDir()

	if err := os.WriteFile(filepath.Join(folder, "window.toml"), []byte("width = 0\nheight = 0\n"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, ok, err := storage.NewWindow(folder).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if ok {
		t.Error("a zero size was reported as recorded")
	}
}
