package opendom

import (
	"errors"
	"testing"

	crowererrors "github.com/fiwon123/crower/internal/errors"
)

func TestOpenCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		path    string
		want    string
		wantErr bool
	}{
		{"windows", "windows", "path.txt", `start "" path.txt`, false},
		{"linux", "linux", "path.txt", "xdg-open 'path.txt'", false},
		{"darwin", "darwin", "path.txt", "open 'path.txt'", false},
		{"unsupported platform", "plan9", "path.txt", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openCommand(tt.goos, tt.path)

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

// A path with a space is quoted for the platform that will run it.
func TestOpenCommandQuotesForThePlatform(t *testing.T) {
	darwin, err := openCommand("darwin", "a b.txt")
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if darwin != "open 'a b.txt'" {
		t.Errorf("darwin command = %q", darwin)
	}

	windows, err := openCommand("windows", "a b.txt")
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if windows != `start "" "a b.txt"` {
		t.Errorf("windows command = %q", windows)
	}
}
