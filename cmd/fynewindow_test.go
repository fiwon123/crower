package cmd

import (
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/desktop"
)

// The window remembers a size that can open, and ignores one that cannot.
func TestFyneWindowRemembersAUsableSize(t *testing.T) {
	window := newFyneWindow(t.TempDir())

	if _, ok := window.Load(); ok {
		t.Fatal("a fresh folder reported a recorded size")
	}

	if err := window.Save(desktop.WindowSize{Width: 1200, Height: 800}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	size, ok := window.Load()
	if !ok || size.Width != 1200 || size.Height != 800 {
		t.Fatalf("Load() = %+v, %v, want 1200x800", size, ok)
	}

	// A window dragged down to nothing must not come back that way.
	if err := window.Save(desktop.WindowSize{Width: 100, Height: 100}); err != nil {
		t.Fatalf("Save unusable: %v", err)
	}

	size, _ = window.Load()
	if size.Width != 1200 || size.Height != 800 {
		t.Errorf("size = %+v, want the last usable one kept", size)
	}
}

// A recorded size too small to work in is not handed to the window.
func TestFyneWindowLoadRejectsAnUnusableRecordedSize(t *testing.T) {
	folder := t.TempDir()

	if err := storage.NewWindow(folder).Save(storage.WindowState{Width: 100, Height: 100}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, ok := newFyneWindow(folder).Load(); ok {
		t.Error("an unusable recorded size was accepted")
	}
}
