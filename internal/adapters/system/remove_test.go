package system_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
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

			got, err := system.RemoveCommand(tt.goos, path, tt.folder)

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
				t.Errorf("RemoveCommand() = %q, want %q", got, tt.want)
			}
		})
	}
}

// The command table only proves the line is right; these two prove the line
// runs. A path with a space is the case that broke before #129, and a folder
// needs a different flag from a file.
func TestRemoverDeletesAFileWithASpaceInItsName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new file.txt")

	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}

	if err := system.NewRemover().File(path); err != nil {
		t.Fatalf("File(%s): %v", path, err)
	}

	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Stat(%s) = %v, want the file gone", path, err)
	}
}

func TestRemoverDeletesAFolderWithASpaceInItsName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new folder")

	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}

	if err := os.WriteFile(filepath.Join(path, "inside.txt"), []byte("x"), 0o600); err != nil {
		t.Fatalf("write inside file: %v", err)
	}

	if err := system.NewRemover().Folder(path); err != nil {
		t.Fatalf("Folder(%s): %v", path, err)
	}

	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Stat(%s) = %v, want the folder and its contents gone", path, err)
	}
}
