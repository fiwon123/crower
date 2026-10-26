package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// addToPathScript is the script under test, relative to this package.
const addToPathScript = "scripts/add_to_path.sh"

// stageScript copies the script into a fresh temporary directory and returns
// the path to the copy. The directory name is a parameter so a test can
// exercise a path containing spaces and parentheses — the reporter's extraction
// directory was "crower_v0.4.0_linux (3)".
func stageScript(t *testing.T, dirName string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("add_to_path.sh is a POSIX shell script")
	}

	data, err := os.ReadFile(addToPathScript)
	if err != nil {
		t.Fatalf("read %s: %v", addToPathScript, err)
	}

	dir := filepath.Join(t.TempDir(), dirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create %s: %v", dir, err)
	}
	path := filepath.Join(dir, "add_to_path.sh")
	if err := os.WriteFile(path, data, 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	return path
}

// runScript runs the script with a throwaway HOME and the given SHELL. The
// real home directory is never touched.
func runScript(t *testing.T, script, home, shell string, prependPath ...string) string {
	t.Helper()

	cmd := exec.Command("sh", script)
	cmd.Dir = filepath.Dir(script)
	cmd.Env = []string{
		"HOME=" + home,
		"SHELL=" + shell,
		"PATH=" + strings.Join(append(prependPath, os.Getenv("PATH")), string(os.PathListSeparator)),
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run %s: %v\n%s", script, err, out)
	}
	return string(out)
}

func readConfig(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// pathLine is the line the script is expected to write for a POSIX shell.
func pathLine(dir string) string {
	return `export PATH="` + dir + `:$PATH"`
}

// TestAddToPathIgnoresOtherShellsConfig is the regression from issue #97: a
// stale entry in ~/.bashrc made the script report success on a zsh machine and
// write nothing to ~/.zshrc.
func TestAddToPathIgnoresOtherShellsConfig(t *testing.T) {
	script := stageScript(t, "crower_v0.4.0_linux (3)")
	dir := filepath.Dir(script)
	home := t.TempDir()

	bashrc := filepath.Join(home, ".bashrc")
	stale := pathLine(dir) + "\n"
	if err := os.WriteFile(bashrc, []byte(stale), 0o644); err != nil {
		t.Fatalf("prepare .bashrc: %v", err)
	}

	out := runScript(t, script, home, "/usr/bin/zsh")

	zshrc, err := os.ReadFile(filepath.Join(home, ".zshrc"))
	if err != nil {
		t.Fatalf("zshrc was not written: %v\noutput:\n%s", err, out)
	}
	if got, want := string(zshrc), pathLine(dir)+"\n"; got != want {
		t.Errorf(".zshrc = %q, want %q", got, want)
	}
	if !strings.Contains(out, ".bashrc") {
		t.Errorf("output does not point at the ignored .bashrc entry:\n%s", out)
	}
	if !strings.Contains(out, "source ~/.zshrc") {
		t.Errorf("output does not tell the user how to apply it:\n%s", out)
	}
	if got := readConfig(t, bashrc); got != stale {
		t.Errorf(".bashrc was modified: %q", got)
	}
}

// TestAddToPathWritesQuotedLine checks that a directory with spaces and
// parentheses survives into the config line intact.
func TestAddToPathWritesQuotedLine(t *testing.T) {
	script := stageScript(t, "crower_v0.4.0_linux (3)")
	dir := filepath.Dir(script)
	home := t.TempDir()

	runScript(t, script, home, "/usr/bin/bash")

	got := readConfig(t, filepath.Join(home, ".bashrc"))
	want := pathLine(dir) + "\n"
	if got != want {
		t.Errorf(".bashrc = %q, want %q", got, want)
	}
}

// TestAddToPathIsIdempotent runs the script twice: the second run must not
// append a duplicate, and must still tell the user how to apply it.
func TestAddToPathIsIdempotent(t *testing.T) {
	script := stageScript(t, "crower")
	dir := filepath.Dir(script)
	home := t.TempDir()

	first := runScript(t, script, home, "/usr/bin/zsh")
	if !strings.Contains(first, "source ~/.zshrc") {
		t.Errorf("first run does not print the source command:\n%s", first)
	}

	second := runScript(t, script, home, "/usr/bin/zsh")
	if !strings.Contains(second, "already configured") {
		t.Errorf("second run does not report the existing entry:\n%s", second)
	}
	if !strings.Contains(second, "source ~/.zshrc") {
		t.Errorf("second run does not print the source command:\n%s", second)
	}

	got := readConfig(t, filepath.Join(home, ".zshrc"))
	if lines := strings.Count(got, "\n"); lines != 1 {
		t.Errorf(".zshrc has %d lines, want 1:\n%s", lines, got)
	}
	if !strings.Contains(got, pathLine(dir)) {
		t.Errorf(".zshrc does not contain %q:\n%s", pathLine(dir), got)
	}
}

// TestAddToPathReportsActiveSession covers the case where the directory is
// already exported: the script should say so instead of asking for a source.
func TestAddToPathReportsActiveSession(t *testing.T) {
	script := stageScript(t, "crower")
	dir := filepath.Dir(script)
	home := t.TempDir()

	out := runScript(t, script, home, "/usr/bin/zsh", dir)

	if !strings.Contains(out, "active in this shell") {
		t.Errorf("output does not report the active session:\n%s", out)
	}
	if strings.Contains(out, "source ~/.zshrc") {
		t.Errorf("output asks for a source that is not needed:\n%s", out)
	}
}

// TestAddToPathFishConfig checks the fish syntax and the nested config
// directory, which does not exist in a fresh HOME.
func TestAddToPathFishConfig(t *testing.T) {
	script := stageScript(t, "crower")
	dir := filepath.Dir(script)
	home := t.TempDir()

	runScript(t, script, home, "/usr/bin/fish")

	got := readConfig(t, filepath.Join(home, ".config", "fish", "config.fish"))
	want := `set -gx PATH "` + dir + `" $PATH` + "\n"
	if got != want {
		t.Errorf("config.fish = %q, want %q", got, want)
	}
}

// TestAddToPathUnknownShellUsesProfile covers the fallback branch.
func TestAddToPathUnknownShellUsesProfile(t *testing.T) {
	script := stageScript(t, "crower")
	dir := filepath.Dir(script)
	home := t.TempDir()

	out := runScript(t, script, home, "/bin/sh")

	got := readConfig(t, filepath.Join(home, ".profile"))
	want := pathLine(dir) + "\n"
	if got != want {
		t.Errorf(".profile = %q, want %q", got, want)
	}
	if !strings.Contains(out, ". ~/.profile") {
		t.Errorf("output does not print the source command:\n%s", out)
	}
}
