package desktop

import "testing"

// A remembered size is only usable when it opens a window someone can work in.
func TestWindowSizeUsable(t *testing.T) {
	tests := []struct {
		name string
		size WindowSize
		want bool
	}{
		{"default", DefaultWindowSize(), true},
		{"minimum", WindowSize{Width: 640, Height: 480}, true},
		{"too narrow", WindowSize{Width: 639, Height: 600}, false},
		{"too short", WindowSize{Width: 800, Height: 479}, false},
		{"zero", WindowSize{}, false},
		{"negative", WindowSize{Width: -800, Height: -600}, false},
		{"absurd", WindowSize{Width: 10001, Height: 800}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.size.Usable(); got != tt.want {
				t.Errorf("Usable() = %v, want %v for %+v", got, tt.want, tt.size)
			}
		})
	}
}
