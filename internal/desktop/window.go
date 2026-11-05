package desktop

// WindowSize is the size a window remembers between runs, in pixels.
//
// It lives here, beside the model, because the command layer records it and the
// window reads it; neither package should have to own the other's type.
type WindowSize struct {
	Width  int
	Height int
}

const (
	minWindowWidth  = 640
	minWindowHeight = 480
	maxWindowWidth  = 10000
	maxWindowHeight = 10000
)

// Usable reports whether a remembered size can open a window someone can work
// in. A corrupted or tiny value falls back to the default rather than opening an
// unusable window that the user cannot resize back.
func (w WindowSize) Usable() bool {
	return w.Width >= minWindowWidth && w.Width <= maxWindowWidth &&
		w.Height >= minWindowHeight && w.Height <= maxWindowHeight
}

// DefaultWindowSize is what the window opens with when nothing usable was
// recorded.
func DefaultWindowSize() WindowSize {
	return WindowSize{Width: 960, Height: 640}
}
