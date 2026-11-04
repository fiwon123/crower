package system_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
	crerrors "github.com/fiwon123/crower/internal/errors"
)

// The runner validates the target before starting anything: a typo in Shell
// must not fall back to a shell the command did not ask for.
func TestRunRejectsAnUnknownShell(t *testing.T) {
	_, err := system.NewRunner().Run("echo hi", "fish", "")
	if !errors.Is(err, crerrors.ErrUnknownShell) {
		t.Fatalf("error = %v, want ErrUnknownShell", err)
	}
}

// The same for Cwd: a folder that is not there is an error, not a chdir the
// shell reports after half the command has run.
func TestRunRejectsAMissingWorkDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")

	_, err := system.NewRunner().Run("echo hi", "", missing)
	if !errors.Is(err, crerrors.ErrInvalidWorkDir) {
		t.Fatalf("error = %v, want ErrInvalidWorkDir", err)
	}
}

// A declared shell runs the line: "sh" works on every platform that has it,
// which is the point of naming a shell instead of guessing one.
func TestRunThroughNamedSh(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test needs a POSIX sh")
	}

	out, err := system.NewRunner().Run("echo hi", "sh", "")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got := strings.TrimSpace(out); got != "hi" {
		t.Errorf("output = %q, want hi", got)
	}
}

// Cwd is the directory the child starts in, not just a field on the struct.
func TestRunUsesTheWorkingDirectory(t *testing.T) {
	dir := t.TempDir()

	out, err := system.NewRunner().Run(printWorkingDirectory(), "", dir)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	assertWorkingDirectory(t, dir, out)
}

// A nested crower step honours the same target: the port changed for both
// entry points, and this is the one easiest to miss.
func TestRunNestedUsesTheWorkingDirectory(t *testing.T) {
	dir := t.TempDir()

	out, err := system.NewRunner().RunNested(printWorkingDirectory(), "", dir)
	if err != nil {
		t.Fatalf("RunNested: %v", err)
	}

	assertWorkingDirectory(t, dir, out)
}

// printWorkingDirectory is the line that reports the child's cwd on each
// platform.
func printWorkingDirectory() string {
	if runtime.GOOS == "windows" {
		return "cd"
	}

	return "pwd"
}

// assertWorkingDirectory compares the directory itself rather than the text:
// Windows case and macOS symlinks make a string comparison flaky.
func assertWorkingDirectory(t *testing.T, want, output string) {
	t.Helper()

	wantInfo, err := os.Stat(want)
	if err != nil {
		t.Fatalf("stat %q: %v", want, err)
	}

	gotInfo, err := os.Stat(strings.TrimSpace(output))
	if err != nil {
		t.Fatalf("stat reported cwd %q: %v", strings.TrimSpace(output), err)
	}

	if !os.SameFile(wantInfo, gotInfo) {
		t.Errorf("working directory = %q, want %q", strings.TrimSpace(output), want)
	}
}
