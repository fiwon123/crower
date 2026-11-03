package createdom

import (
	"errors"
	"testing"

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
			got, err := createFileCommand(tt.goos, "new.txt")

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
			got, err := createFolderCommand(tt.goos, "new")

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

func TestScriptExtension(t *testing.T) {
	tests := []struct {
		goos    string
		want    string
		wantErr bool
	}{
		{"windows", ".bat", false},
		{"linux", ".sh", false},
		{"darwin", ".sh", false},
		{"plan9", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			got, err := scriptExtension(tt.goos)

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
				t.Errorf("extension = %q, want %q", got, tt.want)
			}
		})
	}
}
