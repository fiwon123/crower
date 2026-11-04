package system_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
	crowererrors "github.com/fiwon123/crower/internal/errors"
)

func TestCreateFileCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		want    string
		wantErr bool
	}{
		{"windows", "windows", "type nul > new.txt", false},
		{"linux", "linux", "touch 'new.txt'", false},
		{"darwin", "darwin", "touch 'new.txt'", false},
		{"unsupported platform", "plan9", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := system.CreateFileCommand(tt.goos, "new.txt")

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

func TestCreateFolderCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		want    string
		wantErr bool
	}{
		{"windows", "windows", "mkdir new", false},
		{"linux", "linux", "mkdir -p 'new'", false},
		{"darwin", "darwin", "mkdir -p 'new'", false},
		{"unsupported platform", "plan9", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := system.CreateFolderCommand(tt.goos, "new")

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

// The table proves the line is right; this proves it runs, and that a path
// with a space survives the quoting.
func TestCreatorCreatesAFileAndAFolderWithSpaces(t *testing.T) {
	dir := t.TempDir()
	creator := system.NewCreator()

	file := filepath.Join(dir, "new file.txt")
	if err := creator.File(file); err != nil {
		t.Fatalf("File(%s): %v", file, err)
	}

	if _, err := os.Stat(file); err != nil {
		t.Errorf("Stat(%s) = %v, want the file created", file, err)
	}

	folder := filepath.Join(dir, "new folder", "nested")
	if err := creator.Folder(folder); err != nil {
		t.Fatalf("Folder(%s): %v", folder, err)
	}

	if info, err := os.Stat(folder); err != nil || !info.IsDir() {
		t.Errorf("Stat(%s) = %v, %v; want a directory", folder, info, err)
	}
}
