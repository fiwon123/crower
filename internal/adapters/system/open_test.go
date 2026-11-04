package system

import (
	stderrors "errors"
	"path/filepath"
	"runtime"
	"strings"
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
		{"windows", "windows", "notes.txt", `start "" notes.txt`, false},
		{"windows with a space", "windows", `C:\My Files`, `start "" "C:\My Files"`, false},
		{"linux", "linux", "/tmp/notes.txt", `xdg-open '/tmp/notes.txt'`, false},
		{"linux with a space", "linux", "/tmp/my file", `xdg-open '/tmp/my file'`, false},
		{"darwin", "darwin", "/tmp/notes.txt", `open '/tmp/notes.txt'`, false},
		{"unsupported", "plan9", "notes.txt", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenCommand(tt.goos, tt.path)

			if tt.wantErr {
				if !stderrors.Is(err, crowererrors.ErrUnsupportedOS) {
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

// The opener hands back both the command and its output: one is the diagnostic,
// the other is data a caller may pipe.
func TestOpenerReturnsTheCommandAndOutput(t *testing.T) {
	opener := NewOpener()
	opener.goos = "linux"
	opener.run = func(command string) (string, error) {
		if !strings.Contains(command, "notes.txt") {
			t.Errorf("command = %q, want it to open the path", command)
		}

		return "  opened\n", nil
	}

	command, out, err := opener.Open("notes.txt")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	if command != `xdg-open 'notes.txt'` {
		t.Errorf("command = %q, want the built line", command)
	}

	if out != "opened" {
		t.Errorf("output = %q, want it trimmed", out)
	}
}

func TestOpenerWrapsThePlatformError(t *testing.T) {
	boom := stderrors.New("boom")

	opener := NewOpener()
	opener.goos = "linux"
	opener.run = func(string) (string, error) {
		return "no method available", boom
	}

	command, out, err := opener.Open("notes.txt")
	if !stderrors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}

	if command == "" || out != "no method available" {
		t.Errorf("command = %q, output = %q; want both kept for the message", command, out)
	}

	for _, want := range []string{"open notes.txt", "no method available"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error = %v, want it to contain %q", err, want)
		}
	}
}

func TestOpenerOnAnUnsupportedPlatform(t *testing.T) {
	opener := NewOpener()
	opener.goos = "plan9"
	opener.run = func(string) (string, error) {
		t.Error("the runner was called on an unsupported platform")
		return "", nil
	}

	if _, _, err := opener.Open("notes.txt"); !stderrors.Is(err, crowererrors.ErrUnsupportedOS) {
		t.Fatalf("error = %v, want ErrUnsupportedOS", err)
	}
}

func TestSystemEditorOnWindows(t *testing.T) {
	editor := NewSystemEditor()
	editor.goos = "windows"
	editor.start = func(command string) error {
		if command != "sysdm.cpl" {
			t.Errorf("command = %q, want sysdm.cpl", command)
		}

		return nil
	}
	editor.edit = func(string) error {
		t.Error("the file editor was used on Windows")
		return nil
	}

	if err := editor.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
}

func TestSystemEditorEditsTheShellFile(t *testing.T) {
	tests := []struct {
		goos string
		file string
	}{
		{"linux", ".bashrc"},
		{"darwin", ".zshrc"},
	}

	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			home := t.TempDir()

			editor := NewSystemEditor()
			editor.goos = tt.goos
			editor.home = func() (string, error) { return home, nil }

			var edited string
			editor.edit = func(file string) error {
				edited = file
				return nil
			}

			if err := editor.Open(); err != nil {
				t.Fatalf("Open: %v", err)
			}

			if want := filepath.Join(home, tt.file); edited != want {
				t.Errorf("edited %q, want %q", edited, want)
			}
		})
	}
}

func TestSystemEditorOnAnUnsupportedPlatform(t *testing.T) {
	editor := NewSystemEditor()
	editor.goos = "plan9"

	err := editor.Open()
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("error = %v, want an unsupported-platform error", err)
	}
}

func TestEditorNamePrefersTheEnvironment(t *testing.T) {
	t.Setenv("EDITOR", "my-editor")

	if got := EditorName(); got != "my-editor" {
		t.Errorf("EditorName() = %q, want my-editor", got)
	}
}

func TestShellProfileFileName(t *testing.T) {
	if got := ShellProfileFileName("darwin"); got != ".zshrc" {
		t.Errorf("darwin = %q, want .zshrc", got)
	}

	for _, goos := range []string{"linux", "windows", ""} {
		if got := ShellProfileFileName(goos); got != ".bashrc" {
			t.Errorf("%q = %q, want .bashrc", goos, got)
		}
	}
}

// With no editor named and none on the path there is nothing to run, and the
// error says how to fix it instead of opening something arbitrary.
func TestEditorNameWithNothingAvailable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the PATH lookup is POSIX")
	}

	t.Setenv("EDITOR", "")
	t.Setenv("PATH", t.TempDir())

	if got := EditorName(); got != "" {
		t.Errorf("EditorName() = %q, want empty", got)
	}
}
