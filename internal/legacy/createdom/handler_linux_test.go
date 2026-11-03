//go:build linux

package createdom

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The path has to reach both a login shell and an interactive one: .profile
// is what a login shell reads, .bashrc is what an interactive shell reads,
// and a system that has no .profile at all reads neither.
func TestCreateSystemPathVariableWritesBothFiles(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	handler := Handler{}

	out, err := handler.CreateSystemPathVariable("/opt/mytools")
	if err != nil {
		t.Fatalf("error creating the path variable: %v", err)
	}

	for _, name := range []string{".profile", ".bashrc"} {
		content, err := os.ReadFile(filepath.Join(home, name))
		if err != nil {
			t.Fatalf("error reading %s: %v", name, err)
		}

		if !strings.Contains(string(content), "/opt/mytools") {
			t.Errorf("expected %s to hold the path, got %q", name, string(content))
		}
	}

	if !strings.Contains(out, "new terminal") {
		t.Errorf("expected the message to say what makes the change take effect, got %q", out)
	}
}

// A path that is already there is not added again, and crower says so
// instead of claiming it added it.
func TestCreateSystemPathVariableReportsAnExistingPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	handler := Handler{}

	if _, err := handler.CreateSystemPathVariable("/opt/mytools"); err != nil {
		t.Fatalf("error creating the path variable: %v", err)
	}

	out, err := handler.CreateSystemPathVariable("/opt/mytools")
	if err != nil {
		t.Fatalf("error creating the path variable: %v", err)
	}

	if !strings.Contains(out, "Already in PATH") {
		t.Errorf("expected the message to say the path was already there, got %q", out)
	}

	if strings.Contains(out, "Added to PATH") {
		t.Errorf("expected no claim that the path was added, got %q", out)
	}
}

// The check is a path element check, not a substring one: a path that is the
// prefix of another one is a different path.
func TestCreateSystemPathVariableDoesNotConfuseAPrefixForThePath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	profile := filepath.Join(home, ".profile")
	if err := os.WriteFile(profile, []byte("export PATH=\"$PATH:/opt/tools-extra\"\n"), 0644); err != nil {
		t.Fatalf("error writing the profile: %v", err)
	}

	handler := Handler{}

	if _, err := handler.CreateSystemPathVariable("/opt/tools"); err != nil {
		t.Fatalf("error creating the path variable: %v", err)
	}

	content, err := os.ReadFile(profile)
	if err != nil {
		t.Fatalf("error reading the profile: %v", err)
	}

	if !pathHasValue(string(content), "/opt/tools") {
		t.Errorf("expected the path to be added, got %q", string(content))
	}

	if !pathHasValue(string(content), "/opt/tools-extra") {
		t.Errorf("expected the path that was there to be kept, got %q", string(content))
	}
}

// The PATH the file already declares is kept: the new directory is added to
// it, not in place of it.
func TestCreateSystemPathVariableKeepsThePathTheFileDeclares(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	profile := filepath.Join(home, ".profile")
	if err := os.WriteFile(profile, []byte("export PATH=/usr/local/bin:$PATH\n"), 0644); err != nil {
		t.Fatalf("error writing the profile: %v", err)
	}

	handler := Handler{}

	if _, err := handler.CreateSystemPathVariable("/opt/mytools"); err != nil {
		t.Fatalf("error creating the path variable: %v", err)
	}

	content, err := os.ReadFile(profile)
	if err != nil {
		t.Fatalf("error reading the profile: %v", err)
	}

	if !strings.Contains(string(content), "/usr/local/bin") {
		t.Errorf("expected the declared path to be kept, got %q", string(content))
	}

	if !strings.Contains(string(content), "/opt/mytools") {
		t.Errorf("expected the path to be added, got %q", string(content))
	}
}

// A line that is not a PATH export is not one: "export PATHDIR=…" starts
// with "export PATH" and holds a different variable.
func TestIsPathExport(t *testing.T) {
	cases := []struct {
		line string
		want bool
	}{
		{"export PATH=\"$PATH:/opt\"", true},
		{"export PATH=/usr/local/bin:$PATH", true},
		{"  export PATH=$PATH", true},
		{"export PATH", true},
		{"export PATHDIR=/opt", false},
		{"export MYPATH=/opt", false},
		{"PATH=/opt", false},
		{"# export PATH=/opt", false},
	}

	for _, tc := range cases {
		if got := isPathExport(tc.line); got != tc.want {
			t.Errorf("isPathExport(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

// The elements of the PATH are compared whole.
func TestPathHasValue(t *testing.T) {
	cases := []struct {
		line  string
		value string
		want  bool
	}{
		{"export PATH=\"$PATH:/opt/tools\"", "/opt/tools", true},
		{"export PATH=\"$PATH:/opt/tools-extra\"", "/opt/tools", false},
		{"export PATH=/opt/tools:$PATH", "/opt/tools", true},
		{"export PATH=/opt/tools-extra:$PATH", "/opt/tools", false},
		{"export PATH=\"$PATH:/opt/tools\"", "/opt", false},
	}

	for _, tc := range cases {
		if got := pathHasValue(tc.line, tc.value); got != tc.want {
			t.Errorf("pathHasValue(%q, %q) = %v, want %v", tc.line, tc.value, got, tc.want)
		}
	}
}
