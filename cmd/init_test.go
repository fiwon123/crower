package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShellName(t *testing.T) {
	tests := map[string]string{
		"/bin/bash":           "bash",
		"/usr/bin/zsh":        "zsh",
		"/usr/local/bin/fish": "fish",
		"/opt/pwsh":           "powershell",
		"C:\\Program Files\\PowerShell\\pwsh.exe": "powershell",
		"/bin/sh": "",
		"":        "",
	}

	for shell, want := range tests {
		if got := shellName(shell); got != want {
			t.Errorf("shellName(%q) = %q, want %q", shell, got, want)
		}
	}
}

func TestInitTarget(t *testing.T) {
	home := "/home/u"

	tests := []struct {
		shell string
		file  string
		line  string
	}{
		{"bash", filepath.Join(home, ".bashrc"), `export PATH="$HOME/.local/bin:$PATH"`},
		{"zsh", filepath.Join(home, ".zshrc"), `export PATH="$HOME/.local/bin:$PATH"`},
		{"fish", filepath.Join(home, ".config", "fish", "config.fish"), `fish_add_path "$HOME/.local/bin"`},
		{"powershell", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			file, line := initTarget(tt.shell, home)
			if file != tt.file || line != tt.line {
				t.Errorf("initTarget(%q) = %q, %q; want %q, %q", tt.shell, file, line, tt.file, tt.line)
			}
		})
	}
}

func TestApplyInitIsIdempotent(t *testing.T) {
	home := t.TempDir()
	rc := filepath.Join(home, ".bashrc")
	line := `export PATH="$HOME/.local/bin:$PATH"`

	changed, err := applyInit(rc, line)
	if err != nil || !changed {
		t.Fatalf("first apply = %v, %v; want changed", changed, err)
	}

	changed, err = applyInit(rc, line)
	if err != nil || changed {
		t.Fatalf("second apply = %v, %v; want unchanged", changed, err)
	}

	body, err := os.ReadFile(rc)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	if strings.Count(string(body), line) != 1 {
		t.Errorf("rc = %q, want the line exactly once", body)
	}
}

// A missing startup file is created, including its folder.
func TestApplyInitCreatesTheFile(t *testing.T) {
	rc := filepath.Join(t.TempDir(), ".config", "fish", "config.fish")

	changed, err := applyInit(rc, "fish_add_path")
	if err != nil || !changed {
		t.Fatalf("apply = %v, %v", changed, err)
	}

	if _, err := os.Stat(rc); err != nil {
		t.Errorf("startup file: %v", err)
	}
}
