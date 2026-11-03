package movedom

import (
	"errors"
	"testing"

	crowererrors "github.com/fiwon123/crower/internal/errors"
)

func TestMoveCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		want    string
		wantErr bool
	}{
		{"windows", "windows", "move src.txt dst", false},
		{"linux", "linux", "mv 'src.txt' 'dst'", false},
		{"darwin", "darwin", "mv 'src.txt' 'dst'", false},
		{"unsupported platform", "plan9", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := moveCommand(tt.goos, "src.txt", "dst")

			if tt.wantErr {
				if !errors.Is(err, crowererrors.ErrUnsupportedOS) {
					t.Fatalf("error = %v, want ErrUnsupportedOS", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("error = %v", err)
			}

			if got != tt.want {
				t.Errorf("command = %q, want %q", got, tt.want)
			}
		})
	}
}
