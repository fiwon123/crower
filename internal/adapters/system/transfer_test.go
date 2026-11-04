package system_test

import (
	"errors"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
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
			got, err := system.CopyCommand(tt.goos, "src.txt", "dst", tt.folder)

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
				t.Errorf("CopyCommand() = %q, want %q", got, tt.want)
			}
		})
	}
}

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
			got, err := system.MoveCommand(tt.goos, "src.txt", "dst")

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
				t.Errorf("MoveCommand() = %q, want %q", got, tt.want)
			}
		})
	}
}

// A path with a space or a quote is the reason the builder quotes at all; the
// platform-specific rules are pinned in internal/shell.
func TestCopyCommandQuotesThePaths(t *testing.T) {
	got, err := system.CopyCommand("linux", "a b.txt", "out dir", false)
	if err != nil {
		t.Fatalf("CopyCommand: %v", err)
	}

	if got != "cp 'a b.txt' 'out dir'" {
		t.Errorf("CopyCommand() = %q", got)
	}
}
