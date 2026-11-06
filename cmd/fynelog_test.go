//go:build fyne

package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

// A run that fails leaves fyne.log with what Fyne said.
func TestWindowLogKeepsAFailedRun(t *testing.T) {
	dir := t.TempDir()
	entry := startWindowLog(dir)

	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	log.Printf("Fyne error: failed to initialise GLFW")

	entry.keep()

	body, err := os.ReadFile(filepath.Join(dir, "fyne.log"))
	if err != nil {
		t.Fatalf("fyne.log was not kept: %v", err)
	}

	if !strings.Contains(string(body), "failed to initialise GLFW") {
		t.Errorf("log = %q, want Fyne's reason", body)
	}

	if !strings.Contains(string(body), "crower-gui") {
		t.Errorf("log = %q, want the run to identify itself", body)
	}
}

// A clean run leaves no log behind: the file exists to explain a failure.
func TestWindowLogCleansACleanRun(t *testing.T) {
	dir := t.TempDir()
	entry := startWindowLog(dir)

	entry.clean()

	if _, err := os.Stat(filepath.Join(dir, "fyne.log")); !os.IsNotExist(err) {
		t.Errorf("fyne.log survived a clean run (stat err = %v)", err)
	}
}

// The first Fyne error raises one notification, not one per line.
func TestFirstFyneErrorRaisesOneDialog(t *testing.T) {
	var calls atomic.Int32

	previous := notifyStartupError
	notifyStartupError = func(string) { calls.Add(1) }

	t.Cleanup(func() {
		notifyStartupError = previous
		log.SetOutput(os.Stderr)
	})

	entry := startWindowLog(t.TempDir())

	log.Printf("Fyne error: failed to create window")
	log.Printf("Fyne error: another failure")

	entry.keep()

	if got := calls.Load(); got != 1 {
		t.Errorf("notifications = %d, want exactly one", got)
	}
}

// A line that is not a Fyne error does not raise a dialog.
func TestOrdinaryLogLinesRaiseNoDialog(t *testing.T) {
	var calls atomic.Int32

	previous := notifyStartupError
	notifyStartupError = func(string) { calls.Add(1) }

	t.Cleanup(func() {
		notifyStartupError = previous
		log.SetOutput(os.Stderr)
	})

	entry := startWindowLog(t.TempDir())

	log.Printf("something worth recording")

	entry.keep()

	if got := calls.Load(); got != 0 {
		t.Errorf("notifications = %d, want none", got)
	}
}
