//go:build linux

package system_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
)

func TestSetVariableAppendsTheExport(t *testing.T) {
	dir := home(t)

	out, err := system.NewEnvironment().SetVariable("FOO", "bar")
	if err != nil {
		t.Fatalf("SetVariable: %v", err)
	}

	if !strings.Contains(out, ".bashrc") {
		t.Errorf("output = %q, want it to name the file", out)
	}

	content := readFile(t, filepath.Join(dir, ".bashrc"))
	if !strings.Contains(content, `export FOO="bar"`) {
		t.Errorf("bashrc = %q, want the export line", content)
	}
}

// Adding the same variable twice is refused instead of appending a second
// assignment that would silently win.
func TestSetVariableRefusesADuplicate(t *testing.T) {
	dir := home(t)
	writeFile(t, filepath.Join(dir, ".bashrc"), "export FOO=\"one\"\n")

	_, err := system.NewEnvironment().SetVariable("FOO", "two")
	if err == nil || !strings.Contains(err.Error(), "variable already added") {
		t.Fatalf("error = %v, want variable already added", err)
	}

	content := readFile(t, filepath.Join(dir, ".bashrc"))
	if strings.Contains(content, "two") {
		t.Errorf("bashrc = %q, want it unchanged", content)
	}
}

// The path has to reach both a login shell and an interactive one: .profile
// is what a login shell reads, .bashrc is what an interactive shell reads,
// and a system that has no .profile at all reads neither.
func TestAddPathWritesBothFiles(t *testing.T) {
	dir := home(t)

	out, err := system.NewEnvironment().AddPath("/opt/mytools")
	if err != nil {
		t.Fatalf("AddPath: %v", err)
	}

	for _, name := range []string{".profile", ".bashrc"} {
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}

		if !strings.Contains(string(content), "/opt/mytools") {
			t.Errorf("%s = %q, want it to hold the path", name, content)
		}
	}

	if !strings.Contains(out, "new terminal") {
		t.Errorf("output = %q, want it to say what makes the change take effect", out)
	}
}

// A path that is already there is not added again, and crower says so
// instead of claiming it added it.
func TestAddPathReportsAnExistingPath(t *testing.T) {
	home(t)

	environment := system.NewEnvironment()

	if _, err := environment.AddPath("/opt/mytools"); err != nil {
		t.Fatalf("AddPath: %v", err)
	}

	out, err := environment.AddPath("/opt/mytools")
	if err != nil {
		t.Fatalf("AddPath again: %v", err)
	}

	if !strings.Contains(out, "Already in PATH") {
		t.Errorf("output = %q, want it to say the path was already there", out)
	}

	if strings.Contains(out, "Added to PATH") {
		t.Errorf("output = %q, want no claim that the path was added", out)
	}
}

// The check is a path element check, not a substring one: a path that is the
// prefix of another one is a different path.
func TestAddPathDoesNotConfuseAPrefixForThePath(t *testing.T) {
	dir := home(t)
	profile := filepath.Join(dir, ".profile")

	writeFile(t, profile, "export PATH=\"$PATH:/opt/tools-extra\"\n")

	if _, err := system.NewEnvironment().AddPath("/opt/tools"); err != nil {
		t.Fatalf("AddPath: %v", err)
	}

	content := readFile(t, profile)
	if !strings.Contains(content, "/opt/tools\"") && !strings.Contains(content, "/opt/tools:") {
		t.Errorf("profile = %q, want the path added", content)
	}

	if !strings.Contains(content, "/opt/tools-extra") {
		t.Errorf("profile = %q, want the path that was there kept", content)
	}
}

// The PATH the file already declares is kept: the new directory is added to
// it, not in place of it.
func TestAddPathKeepsThePathTheFileDeclares(t *testing.T) {
	dir := home(t)
	profile := filepath.Join(dir, ".profile")

	writeFile(t, profile, "export PATH=/usr/local/bin:$PATH\n")

	if _, err := system.NewEnvironment().AddPath("/opt/mytools"); err != nil {
		t.Fatalf("AddPath: %v", err)
	}

	content := readFile(t, profile)
	if !strings.Contains(content, "/usr/local/bin") {
		t.Errorf("profile = %q, want the declared path kept", content)
	}

	if !strings.Contains(content, "/opt/mytools") {
		t.Errorf("profile = %q, want the path added", content)
	}
}
