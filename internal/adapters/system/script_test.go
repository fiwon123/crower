package system_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
	crowererrors "github.com/fiwon123/crower/internal/errors"
)

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
			got, err := system.ScriptExtension(tt.goos)

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

func TestScriptsCreateMakesTheFileAndReportsItsFolder(t *testing.T) {
	folder := filepath.Join(t.TempDir(), "scripts")

	path, reported, err := system.NewScripts(folder).Create("deploy")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if reported != folder {
		t.Errorf("folder = %q, want %q", reported, folder)
	}

	if filepath.Dir(path) != folder {
		t.Errorf("path = %q, want it inside %q", path, folder)
	}

	if filepath.Ext(path) != ".sh" && filepath.Ext(path) != ".bat" {
		t.Errorf("path = %q, want a script extension", path)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("Stat(%s) = %v, want the script created", path, err)
	}
}

// Creating the command again must not wipe what the user wrote in the script.
func TestScriptsCreateKeepsAnExistingFile(t *testing.T) {
	folder := filepath.Join(t.TempDir(), "scripts")

	path, _, err := system.NewScripts(folder).Create("deploy")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	body := []byte("#!/bin/sh\necho deploy\n")
	if err := os.WriteFile(path, body, 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	again, _, err := system.NewScripts(folder).Create("deploy")
	if err != nil {
		t.Fatalf("Create again: %v", err)
	}

	content, err := os.ReadFile(again)
	if err != nil {
		t.Fatalf("read script: %v", err)
	}

	if string(content) != string(body) {
		t.Errorf("script = %q, want the existing body kept", content)
	}
}
