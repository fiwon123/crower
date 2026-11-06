package main

import (
	"archive/zip"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckNames(t *testing.T) {
	tests := []struct {
		name    string
		gui     bool
		entries []string
		wantOK  bool
		wantMsg string
	}{
		{
			name:    "plain unix names",
			entries: []string{"crower", "cr", "add_to_path.sh", "README.md", "NOTICE"},
			wantOK:  true,
		},
		{
			name:    "plain windows names",
			entries: []string{"crower.exe", "cr.exe", "add_to_path.bat", "README.md", "NOTICE"},
			wantOK:  true,
		},
		{
			name:    "the v0.4.0 regression",
			entries: []string{"crower_linux_amd64", "cr_linux_amd64", "add_to_path.sh"},
			wantMsg: "missing \"crower\"",
		},
		{
			name:    "the v0.4.0 windows regression",
			entries: []string{"crower_windows_amd64.exe", "cr_windows_amd64.exe"},
			wantMsg: "arch-suffixed binary",
		},
		{
			name:    "short name missing",
			entries: []string{"crower", "add_to_path.sh"},
			wantMsg: "missing \"cr\"",
		},
		{
			name:    "long name missing",
			entries: []string{"cr"},
			wantMsg: "missing \"crower\"",
		},
		{
			name:    "empty archive",
			entries: nil,
			wantMsg: "missing \"crower\"",
		},
		{
			name:    "parent traversal",
			entries: []string{"crower", "cr", "../README.md"},
			wantMsg: "unsafe member",
		},
		{
			name:    "absolute member",
			entries: []string{"crower", "cr", "/etc/passwd"},
			wantMsg: "unsafe member",
		},
		{
			name:    "windows parent traversal",
			entries: []string{"crower.exe", "cr.exe", `..\README.md`},
			wantMsg: "unsafe member",
		},
		{
			name:    "gui archive with all three",
			gui:     true,
			entries: []string{"crower", "cr", "crower-gui", "README.md", "NOTICE"},
			wantOK:  true,
		},
		{
			name:    "gui archive windows",
			gui:     true,
			entries: []string{"crower.exe", "cr.exe", "crower-gui.exe", "README.md", "NOTICE"},
			wantOK:  true,
		},
		{
			name:    "gui archive without the gui binary",
			gui:     true,
			entries: []string{"crower", "cr"},
			wantMsg: "missing \"crower-gui\"",
		},
		{
			name:    "gui archive without the cli pair",
			gui:     true,
			entries: []string{"crower-gui"},
			wantMsg: "missing \"crower\"",
		},
		{
			name:    "archive without the notices",
			entries: []string{"crower", "cr", "README.md"},
			wantMsg: "missing \"NOTICE\"",
		},
		{
			name:    "cli archive may still carry the gui binary",
			entries: []string{"crower", "cr", "crower-gui", "README.md", "NOTICE"},
			wantOK:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problems := checkNames(tt.gui, tt.entries)
			if tt.wantOK {
				if len(problems) != 0 {
					t.Fatalf("checkNames(%v, %v) = %v, want no problems", tt.gui, tt.entries, problems)
				}
				return
			}
			if len(problems) == 0 {
				t.Fatalf("checkNames(%v, %v) passed, want a problem", tt.gui, tt.entries)
			}
			joined := strings.Join(problems, "\n")
			if !strings.Contains(joined, tt.wantMsg) {
				t.Errorf("checkNames(%v, %v) = %q, want it to mention %q", tt.gui, tt.entries, joined, tt.wantMsg)
			}
		})
	}
}

func TestIsUnsafePath(t *testing.T) {
	tests := []struct {
		entry string
		want  bool
	}{
		{"crower", false},
		{"crower-gui", false},
		{"README.md", false},
		{"add_to_path.sh", false},
		{"../README.md", true},
		{"a/../b", true},
		{`..\README.md`, true},
		{"/etc/passwd", true},
		{`C:\Windows\system32`, true},
		{"a..b", false},
	}

	for _, tt := range tests {
		if got := isUnsafePath(tt.entry); got != tt.want {
			t.Errorf("isUnsafePath(%q) = %v, want %v", tt.entry, got, tt.want)
		}
	}
}

func TestIsArchSuffixed(t *testing.T) {
	tests := []struct {
		base string
		want bool
	}{
		{"crower_linux_amd64", true},
		{"cr_linux_amd64", true},
		{"crower_windows_amd64.exe", true},
		{"cr_darwin_arm64", true},
		{"crower", false},
		{"cr", false},
		{"crower-gui", false},
		{"crower.exe", false},
		{"README.md", false},
		{"add_to_path.sh", false},
		{"crowex_linux_amd64", false},
	}

	for _, tt := range tests {
		if got := isArchSuffixed(tt.base); got != tt.want {
			t.Errorf("isArchSuffixed(%q) = %v, want %v", tt.base, got, tt.want)
		}
	}
}

// buildWindows compiles a minimal executable for one Windows architecture.
func buildWindows(t *testing.T, dir, goarch string, gui bool) string {
	t.Helper()

	source := filepath.Join(dir, "fixture_"+goarch)
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(source, "go.mod"), []byte("module fixture\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	if err := os.WriteFile(filepath.Join(source, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	output := filepath.Join(source, "probe.exe")

	args := []string{"build", "-o", output}
	if gui {
		args = append(args, "-ldflags", "-H windowsgui")
	}

	args = append(args, ".")

	command := exec.Command("go", args...)
	command.Dir = source
	command.Env = append(os.Environ(), "GOOS=windows", "GOARCH="+goarch, "CGO_ENABLED=0")

	if out, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build fixture for %s: %v\n%s", goarch, err, out)
	}

	return output
}

// zipArchive writes a zip whose members are copied from the given paths.
func zipArchive(t *testing.T, path string, members map[string]string) {
	t.Helper()

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create archive: %v", err)
	}
	defer file.Close()

	writer := zip.NewWriter(file)

	for name, source := range members {
		body, err := os.ReadFile(source)
		if err != nil {
			t.Fatalf("read %s: %v", source, err)
		}

		entry, err := writer.Create(name)
		if err != nil {
			t.Fatalf("add %s: %v", name, err)
		}

		if _, err := entry.Write(body); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close archive: %v", err)
	}
}

// An archive whose members are built for another machine is rejected: this is
// the "This app can't run on your PC" failure, caught before publication.
func TestCheckWindowsBinariesCatchesTheWrongArchitecture(t *testing.T) {
	dir := t.TempDir()
	amd64 := buildWindows(t, dir, "amd64", false)
	arm64 := buildWindows(t, dir, "arm64", false)

	note := filepath.Join(dir, "WINDOWS.txt")
	if err := os.WriteFile(note, []byte("architecture note"), 0o644); err != nil {
		t.Fatalf("write note: %v", err)
	}

	archive := filepath.Join(dir, "crower_v0.20.16_windows_amd64.zip")
	zipArchive(t, archive, map[string]string{
		"crower.exe":     amd64,
		"cr.exe":         amd64,
		"crower-gui.exe": arm64,
		"WINDOWS.txt":    note,
	})

	problems := checkWindowsBinaries(archive, []string{"crower.exe", "cr.exe", "crower-gui.exe", "WINDOWS.txt"})

	joined := strings.Join(problems, "\n")
	if !strings.Contains(joined, "crower-gui.exe") || !strings.Contains(joined, "arm64") {
		t.Fatalf("problems = %q, want the arm64 window binary rejected", joined)
	}
}

// A window binary without its resources is rejected too, and a Windows archive
// without the architecture note is incomplete.
func TestCheckWindowsBinariesWantsResourcesAndTheNote(t *testing.T) {
	dir := t.TempDir()
	amd64 := buildWindows(t, dir, "amd64", false)
	gui := buildWindows(t, dir, "amd64", true)

	archive := filepath.Join(dir, "crower_v0.20.16_windows_amd64.zip")
	zipArchive(t, archive, map[string]string{
		"crower.exe":     amd64,
		"cr.exe":         amd64,
		"crower-gui.exe": gui,
	})

	problems := checkWindowsBinaries(archive, []string{"crower.exe", "cr.exe", "crower-gui.exe"})

	joined := strings.Join(problems, "\n")
	for _, want := range []string{"WINDOWS.txt", "resource"} {
		if !strings.Contains(joined, want) {
			t.Errorf("problems = %q, want it to mention %q", joined, want)
		}
	}
}
