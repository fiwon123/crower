package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage"
)

// The window shows the version the binary was built from, not a guess from the
// environment.
func TestFyneUpdaterCurrentIsTheBuildVersion(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	if got := newFyneUpdater(cfg, "v1.2.3").Current(); got != "v1.2.3" {
		t.Errorf("Current() = %q, want v1.2.3", got)
	}
}

// Restore previous uses the recorded version list, so a window can offer the
// same rollback the CLI has.
func TestFyneUpdaterReportsTheRecordedPreviousVersion(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	versions := storage.NewVersions(cfg.ConfigFolder())
	if err := versions.RecordInstalled("v1.2.2"); err != nil {
		t.Fatalf("RecordInstalled: %v", err)
	}

	previous, found := newFyneUpdater(cfg, "v1.2.3").Previous()
	if !found || previous != "v1.2.2" {
		t.Fatalf("Previous() = %q, %v; want v1.2.2, true", previous, found)
	}
}

// A fresh installation has nothing to restore, and the window must say so
// rather than offer a button that cannot work.
func TestFyneUpdaterHasNoPreviousVersionOnAFreshInstallation(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	if previous, found := newFyneUpdater(cfg, "v1.2.3").Previous(); found {
		t.Errorf("Previous() = %q, true; want no version", previous)
	}
}

// Restart spawns the replacement exactly once and never runs it in a test.
func TestFyneUpdaterRestartSpawnsTheReplacement(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	spawned := 0
	updater := newFyneUpdater(cfg, "v1.2.3")
	updater.executable = "/opt/crower-gui"
	updater.args = []string{"--flag"}

	var gotExecutable string
	var gotArgs []string

	updater.start = func(executable string, args []string) error {
		spawned++
		gotExecutable = executable
		gotArgs = args

		return nil
	}

	if err := updater.Restart(); err != nil {
		t.Fatalf("Restart: %v", err)
	}

	if spawned != 1 {
		t.Errorf("spawned %d times, want 1", spawned)
	}

	if gotExecutable != "/opt/crower-gui" {
		t.Errorf("executable = %q, want the captured path", gotExecutable)
	}

	if !slices.Equal(gotArgs, []string{"--flag"}) {
		t.Errorf("args = %v, want the captured arguments", gotArgs)
	}
}

// The executable is captured when the updater is built, before an install moves
// the running binary aside.
func TestFyneUpdaterCapturesTheExecutable(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	if updater := newFyneUpdater(cfg, "v1.2.3"); updater.executable == "" {
		t.Error("the updater did not capture an executable path")
	}
}

// A path whose file was already replaced may carry the " (deleted)" suffix on
// Linux; the replacement argv uses the path without it.
func TestReplacementArgvStripsADeletedSuffix(t *testing.T) {
	argv := replacementArgv("/opt/crower-gui (deleted)", []string{"--x"})

	if !slices.Equal(argv, []string{"/opt/crower-gui", "--x"}) {
		t.Errorf("argv = %v, want [/opt/crower-gui --x]", argv)
	}
}

// A failed restart has no window to report to, so it leaves a note next to the
// configuration.
func TestFyneUpdaterLogsARestartFailure(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	updater := newFyneUpdater(cfg, "v1.2.3")
	updater.start = func(string, []string) error { return errors.New("boom") }

	if err := updater.Restart(); err == nil {
		t.Fatal("Restart accepted a failing spawn")
	}

	data, err := os.ReadFile(filepath.Join(cfg.ConfigFolder(), "restart-error.log"))
	if err != nil {
		t.Fatalf("restart-error.log: %v", err)
	}

	if !strings.Contains(string(data), "boom") {
		t.Errorf("log = %q, want the failure", data)
	}
}

// The window installs the window member. update.New() extracts the CLI, and
// selfupdate writes those bytes over crower-gui, so Restart launched the CLI —
// which prints help and exits on Linux (#489, #513). #490 added NewWindow but
// never wired it; this test is the guard that was missing.
func TestFyneUpdaterInstallsTheWindowMember(t *testing.T) {
	cfg := loadFyneConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))

	updater := newFyneUpdater(cfg, "v1.2.3")
	if updater.installer == nil {
		t.Fatal("the window updater has no installer")
	}

	if got := updater.installer.Member(); got != "crower-gui" {
		t.Errorf("installer member = %q, want crower-gui", got)
	}
}
