package copydom

import (
	"errors"
	"testing"

	crowererrors "github.com/fiwon123/crower/internal/errors"
)

func TestCopyCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		folder  bool
		want    string
		wantErr bool
	}{
		{"windows file", "windows", false, "copy src.txt dst", false},
		{"windows folder", "windows", true, "xcopy src.txt dst /E /I", false},
		{"linux file", "linux", false, "cp 'src.txt' 'dst'", false},
		{"linux folder", "linux", true, "cp -r 'src.txt' 'dst'", false},
		{"darwin file", "darwin", false, "cp 'src.txt' 'dst'", false},
		{"darwin folder", "darwin", true, "cp -r 'src.txt' 'dst'", false},
		{"unsupported platform", "plan9", false, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := copyCommand(tt.goos, "src.txt", "dst", tt.folder)

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
func TestCopyCommandQuotesForThePlatform(t *testing.T) {
	posix, err := copyCommand("darwin", "a b.txt", "dst", false)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if posix != "cp 'a b.txt' 'dst'" {
		t.Errorf("darwin command = %q", posix)
	}

	windows, err := copyCommand("windows", "a b.txt", "dst", false)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if windows != `copy "a b.txt" dst` {
		t.Errorf("windows command = %q", windows)
	}
}
