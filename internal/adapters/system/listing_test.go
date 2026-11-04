package system_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
	crowererrors "github.com/fiwon123/crower/internal/errors"
)

// The platform dispatch is tested through the pure command builders rather than
// through the methods: the builders decide what runs, they need no host, and
// they never spawn a shell.
func TestFolderListingCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		path    string
		want    string
		wantErr bool
	}{
		{"windows", "windows", "/tmp/example", `dir /tmp/example`, false},
		{"windows with a space", "windows", `C:\My Files`, `dir "C:\My Files"`, false},
		{"linux", "linux", "/tmp/example", `ls -la -- '/tmp/example'`, false},
		{"darwin", "darwin", "/tmp/example", `ls -la -- '/tmp/example'`, false},
		{"unsupported freebsd", "freebsd", "/tmp/example", "", true},
		{"unsupported empty", "", "/tmp/example", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := system.FolderListingCommand(tt.goos, tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && !errors.Is(err, crowererrors.ErrUnsupportedOS) {
				t.Errorf("error = %v, want ErrUnsupportedOS", err)
			}

			if got != tt.want {
				t.Errorf("command = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVariablesListingCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		want    string
		wantErr bool
	}{
		{"windows", "windows", "set", false},
		{"linux", "linux", "printenv", false},
		{"darwin", "darwin", "printenv", false},
		{"unsupported freebsd", "freebsd", "", true},
		{"unsupported empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := system.VariablesListingCommand(tt.goos)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && !errors.Is(err, crowererrors.ErrUnsupportedOS) {
				t.Errorf("error = %v, want ErrUnsupportedOS", err)
			}

			if got != tt.want {
				t.Errorf("command = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPathListingCommand(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		want    string
		wantErr bool
	}{
		{"windows quotes the value", "windows", `echo "%PATH%"`, false},
		{"linux", "linux", `echo "$PATH"`, false},
		{"darwin", "darwin", `echo "$PATH"`, false},
		{"unsupported freebsd", "freebsd", "", true},
		{"unsupported empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := system.PathListingCommand(tt.goos)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && !errors.Is(err, crowererrors.ErrUnsupportedOS) {
				t.Errorf("error = %v, want ErrUnsupportedOS", err)
			}

			if got != tt.want {
				t.Errorf("command = %q, want %q", got, tt.want)
			}
		})
	}
}

// The structured folder result reads the directory, so it describes entries on
// every platform without parsing a listing.
func TestFolderEntries(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	rows, err := system.NewListings().FolderEntries(dir)
	if err != nil {
		t.Fatalf("FolderEntries: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("got %d entries, want 2: %+v", len(rows), rows)
	}

	byName := map[string]bool{}
	for _, row := range rows {
		byName[row.Name] = row.Dir

		if row.Name == "a.txt" {
			if row.Size != 5 {
				t.Errorf("a.txt size = %d, want 5", row.Size)
			}

			if !strings.Contains(row.Mode, "rw") {
				t.Errorf("a.txt mode = %q, want a readable mode", row.Mode)
			}

			if row.Modified == "" {
				t.Error("a.txt has no modification time")
			}
		}
	}

	if byName["a.txt"] {
		t.Error("a.txt is reported as a directory")
	}

	if !byName["sub"] {
		t.Error("sub is not reported as a directory")
	}
}

func TestFolderEntriesReportsAMissingFolder(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")

	_, err := system.NewListings().FolderEntries(missing)
	if err == nil || !strings.Contains(err.Error(), "list folder") {
		t.Fatalf("error = %v, want list folder", err)
	}
}
