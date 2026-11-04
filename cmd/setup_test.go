package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/session"
)

func TestSetupConfigPath(t *testing.T) {
	home := "/home/u"

	got, err := setupConfigPath(home, "streamer")
	if err != nil {
		t.Fatalf("setupConfigPath: %v", err)
	}

	if want := filepath.Join(home, "crower", "setups", "streamer", "crower.yaml"); got != want {
		t.Errorf("path = %q, want %q", got, want)
	}

	for _, name := range []string{"", ".", "..", "a/b", `a\b`, "a b"} {
		if _, err := setupConfigPath(home, name); err == nil {
			t.Errorf("setupConfigPath(%q) returned no error", name)
		}
	}
}

// The selection order is explicit path, then setup, then the default
// installation.
func TestResolveConfigPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	explicit := filepath.Join(home, "explicit.yaml")

	got, err := resolveConfigPath(explicit, "streamer")
	if err != nil {
		t.Fatalf("resolveConfigPath: %v", err)
	}

	if got != explicit {
		t.Errorf("explicit = %q, want %q", got, explicit)
	}

	got, err = resolveConfigPath("", "streamer")
	if err != nil {
		t.Fatalf("resolveConfigPath: %v", err)
	}

	if want := filepath.Join(home, "crower", "setups", "streamer", "crower.yaml"); got != want {
		t.Errorf("setup = %q, want %q", got, want)
	}

	got, err = resolveConfigPath("", "")
	if err != nil {
		t.Fatalf("resolveConfigPath: %v", err)
	}

	if want := filepath.Join(home, "crower", "crower.yaml"); got != want {
		t.Errorf("default = %q, want %q", got, want)
	}
}

// runSetup executes the root command against a temporary HOME, without the
// --config every other helper passes, so the setup resolution is what decides
// the path.
func runSetup(t *testing.T, home string, args ...string) (string, string, error) {
	t.Helper()

	resetFlags()
	t.Setenv("HOME", home)

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	previous := streams
	streams = session.Streams{In: strings.NewReader(""), Out: out, Err: errOut}
	t.Cleanup(func() { streams = previous })

	rootCmd.SetArgs(args)
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	err := rootCmd.Execute()

	return out.String(), errOut.String(), err
}

func TestSetupLifecycle(t *testing.T) {
	home := t.TempDir()

	if _, _, err := runSetup(t, home, "setup", "new", "streamer"); err != nil {
		t.Fatalf("setup new: %v", err)
	}

	if _, _, err := runSetup(t, home, "setup", "use", "streamer"); err != nil {
		t.Fatalf("setup use: %v", err)
	}

	// The persisted default now points the plain commands at streamer.
	if _, _, err := runSetup(t, home, "create", "obs", "echo obs"); err != nil {
		t.Fatalf("create in streamer: %v", err)
	}

	stdout, _, err := runSetup(t, home, "list")
	if err != nil {
		t.Fatalf("list in streamer: %v", err)
	}

	if !strings.Contains(stdout, "obs") {
		t.Errorf("streamer list = %q, want the command", stdout)
	}

	// Switching back hides it: the setups are isolated.
	if _, _, err := runSetup(t, home, "setup", "use", "default"); err != nil {
		t.Fatalf("setup use default: %v", err)
	}

	stdout, _, err = runSetup(t, home, "list")
	if err != nil {
		t.Fatalf("list in default: %v", err)
	}

	if strings.Contains(stdout, "obs") {
		t.Errorf("default list = %q, want it empty", stdout)
	}

	// show reports the active setup.
	stdout, _, err = runSetup(t, home, "setup", "show")
	if err != nil {
		t.Fatalf("setup show: %v", err)
	}

	if !strings.Contains(stdout, "Setup:    default") {
		t.Errorf("show = %q, want the default setup", stdout)
	}
}

func TestSetupUseRejectsAnUnknownSetup(t *testing.T) {
	home := t.TempDir()

	_, _, err := runSetup(t, home, "setup", "use", "ghost")
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("error = %v, want does-not-exist", err)
	}
}

func TestSetupDeleteProtectsTheActiveSetup(t *testing.T) {
	home := t.TempDir()

	if _, _, err := runSetup(t, home, "setup", "new", "gaming"); err != nil {
		t.Fatalf("setup new: %v", err)
	}

	if _, _, err := runSetup(t, home, "setup", "use", "gaming"); err != nil {
		t.Fatalf("setup use: %v", err)
	}

	_, _, err := runSetup(t, home, "setup", "delete", "gaming")
	if err == nil || !strings.Contains(err.Error(), "is active") {
		t.Fatalf("error = %v, want the active-setup refusal", err)
	}

	if _, _, err := runSetup(t, home, "setup", "delete", "gaming", "--force"); err != nil {
		t.Fatalf("forced delete: %v", err)
	}
}

func TestSetupNewRejectsABadName(t *testing.T) {
	home := t.TempDir()

	_, _, err := runSetup(t, home, "setup", "new", "a/b")
	if err == nil || !strings.Contains(err.Error(), "invalid setup name") {
		t.Fatalf("error = %v, want invalid name", err)
	}
}
