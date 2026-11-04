package cmd

import (
	"path/filepath"
	"testing"
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
