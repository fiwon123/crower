//go:build linux

package system_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
)

// home points HOME at a fresh directory and returns it. The environment file
// the adapter writes is read from HOME, so a test must never touch the real one.
func home(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	t.Setenv("HOME", dir)

	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return string(content)
}

func TestVariableEmptiesTheExportLine(t *testing.T) {
	dir := home(t)
	bashrc := filepath.Join(dir, ".bashrc")

	writeFile(t, bashrc, "export FOO=bar\nexport KEEP=1\n")

	out, err := system.NewEnvironment().Variable("FOO")
	if err != nil {
		t.Fatalf("Variable(FOO): %v", err)
	}

	if out != "var name deleted" {
		t.Errorf("output = %q, want %q", out, "var name deleted")
	}

	// Only the matching line is emptied; everything else is left in place.
	if got, want := readFile(t, bashrc), "\nexport KEEP=1\n"; got != want {
		t.Errorf("bashrc = %q, want %q", got, want)
	}
}

func TestVariableThatIsNotThereIsAnError(t *testing.T) {
	dir := home(t)
	bashrc := filepath.Join(dir, ".bashrc")

	writeFile(t, bashrc, "export KEEP=1\n")

	_, err := system.NewEnvironment().Variable("MISSING")
	if err == nil {
		t.Fatal("Variable(MISSING) returned no error")
	}

	if !strings.Contains(err.Error(), "var name not found") {
		t.Errorf("error = %v, want var name not found", err)
	}

	if got, want := readFile(t, bashrc), "export KEEP=1\n"; got != want {
		t.Errorf("bashrc changed to %q, want %q", got, want)
	}
}

func TestVariablePathRemovesOneEntryAndItsSeparator(t *testing.T) {
	dir := home(t)
	profile := filepath.Join(dir, ".profile")

	writeFile(t, profile, "export PATH=/a:/b:/c\n")

	out, err := system.NewEnvironment().VariablePath("/b")
	if err != nil {
		t.Fatalf("VariablePath(/b): %v", err)
	}

	if out != "Deleted path from PATH" {
		t.Errorf("output = %q, want %q", out, "Deleted path from PATH")
	}

	if got, want := readFile(t, profile), "export PATH=/a:/c\n"; got != want {
		t.Errorf("profile = %q, want %q", got, want)
	}
}

func TestVariablePathThatIsNotThereIsAnError(t *testing.T) {
	dir := home(t)
	profile := filepath.Join(dir, ".profile")

	writeFile(t, profile, "export PATH=/a:/b\n")

	_, err := system.NewEnvironment().VariablePath("/missing")
	if err == nil {
		t.Fatal("VariablePath(/missing) returned no error")
	}

	if !strings.Contains(err.Error(), "path not found") {
		t.Errorf("error = %v, want path not found", err)
	}

	if got, want := readFile(t, profile), "export PATH=/a:/b\n"; got != want {
		t.Errorf("profile changed to %q, want %q", got, want)
	}
}

// A line that mentions PATH but is not the export line ends the search, so a
// test or a comment below it cannot be mistaken for the variable.
func TestVariablePathStopsAtTheFirstOtherPATHLine(t *testing.T) {
	dir := home(t)
	profile := filepath.Join(dir, ".profile")

	writeFile(t, profile, "export PATH_EXTRA=/a:/b\n")

	_, err := system.NewEnvironment().VariablePath("/a")
	if err == nil {
		t.Fatal("VariablePath(/a) returned no error")
	}

	if got, want := readFile(t, profile), "export PATH_EXTRA=/a:/b\n"; got != want {
		t.Errorf("profile changed to %q, want %q", got, want)
	}
}
