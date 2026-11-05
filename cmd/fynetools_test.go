package cmd

import (
	"testing"
)

// The capture keeps what a use case prints, one entry per line, and drops the
// blank lines a command may emit.
func TestCaptureOutputCollectsResultInfoAndEmit(t *testing.T) {
	capture := &captureOutput{}

	capture.Info("opened ", "/tmp/a", " with ", "xdg-open")
	capture.Result("extracted 3 files\n")
	capture.Result("   ")
	if err := capture.Emit("found: a.txt", nil); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	want := "opened /tmp/a with xdg-open\nextracted 3 files\nfound: a.txt"
	if got := capture.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

// An unknown search mode is refused before any adapter runs.
func TestFyneToolsRefusesAnUnknownSearchMode(t *testing.T) {
	cfg := loadFyneConfigAt(t, t.TempDir()+"/crower.yaml")

	if _, _, err := newFyneTools(cfg).Search("nonsense", "needle", "."); err == nil {
		t.Fatal("Search accepted an unknown mode")
	}
}
