//go:build fyne

package fyneui

import (
	"testing"

	"fyne.io/fyne/v2"

	"github.com/fiwon123/crower/internal/desktop"
)

type stubWindowStore struct {
	size  desktop.WindowSize
	ok    bool
	saved []desktop.WindowSize
}

func (s *stubWindowStore) Load() (desktop.WindowSize, bool) { return s.size, s.ok }

func (s *stubWindowStore) Save(size desktop.WindowSize) error {
	s.saved = append(s.saved, size)

	return nil
}

// The window opens at the remembered size when there is a usable one, and at the
// default otherwise.
func TestInitialSize(t *testing.T) {
	defaultSize := fyne.NewSize(960, 640)

	if got := initialSize(nil); got != defaultSize {
		t.Errorf("initialSize(nil) = %v, want %v", got, defaultSize)
	}

	remembered := &stubWindowStore{size: desktop.WindowSize{Width: 1200, Height: 800}, ok: true}

	if got, want := initialSize(remembered), fyne.NewSize(1200, 800); got != want {
		t.Errorf("initialSize = %v, want %v", got, want)
	}

	unusable := &stubWindowStore{size: desktop.WindowSize{Width: 10, Height: 10}, ok: true}

	if got := initialSize(unusable); got != defaultSize {
		t.Errorf("initialSize = %v, want the default for an unusable size", got)
	}

	missing := &stubWindowStore{}

	if got := initialSize(missing); got != defaultSize {
		t.Errorf("initialSize = %v, want the default when nothing was recorded", got)
	}
}

// A canvas size is only worth remembering when it can open a usable window.
func TestRecordedSize(t *testing.T) {
	size, ok := recordedSize(fyne.NewSize(1000, 700))
	if !ok || size.Width != 1000 || size.Height != 700 {
		t.Errorf("recordedSize = %+v, %v, want 1000x700", size, ok)
	}

	if _, ok := recordedSize(fyne.NewSize(10, 10)); ok {
		t.Error("an unusable canvas size was accepted")
	}
}
