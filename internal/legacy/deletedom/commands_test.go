package deletedom

import (
	"errors"
	"testing"

	crowererrors "github.com/fiwon123/crower/internal/errors"
)

func TestRemoveCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		folder  bool
		want    string
		wantErr bool
	}{
		{"windows file", "windows", false, "del path.txt", false},
		{"windows folder", "windows", true, "rmdir /s /q path", false},
		{"linux file", "linux", false, "rm 'path.txt'", false},
		{"linux folder", "linux", true, "rm -r 'path'", false},
		{"darwin file", "darwin", false, "rm 'path.txt'", false},
		{"darwin folder", "darwin", true, "rm -r 'path'", false},
		{"unsupported platform", "plan9", false, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "path.txt"
			if tt.folder {
				path = "path"
			}

			got, err := removeCommand(tt.goos, path, tt.folder)

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
