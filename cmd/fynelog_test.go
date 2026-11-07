//go:build fyne

package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// A run that fails leaves fyne.log with what Fyne said.
func TestWindowLogKeepsAFailedRun(t *testing.T) {
	dir := t.TempDir()
	entry := startWindowLog(dir)

	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	// A recoverable Fyne error: the fatal window-creation path exits the
	// process, which is the right behaviour but not something a unit test can
	// survive.
	log.Printf("Fyne error: the clipboard is unavailable")

	entry.keep()

	body, err := os.ReadFile(filepath.Join(dir, "fyne.log"))
	if err != nil {
		t.Fatalf("fyne.log was not kept: %v", err)
	}

	if !strings.Contains(string(body), "the clipboard is unavailable") {
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

// A recoverable Fyne message stays in the log. A locale Fyne cannot round-trip
// is the common one, and a modal box for it is a startup error the user has to
// dismiss for nothing (#471).
func TestRecoverableFyneErrorRaisesNoDialog(t *testing.T) {
	var calls atomic.Int32

	previous := notifyStartupError
	notifyStartupError = func(string) { calls.Add(1) }

	t.Cleanup(func() {
		notifyStartupError = previous
		log.SetOutput(os.Stderr)
	})

	entry := startWindowLog(t.TempDir())

	log.Printf("Fyne error:  Error parsing user locale pt-BR-Latn")

	entry.keep()

	time.Sleep(500 * time.Millisecond)

	if got := calls.Load(); got != 0 {
		t.Errorf("notifications = %d, want none for a recoverable message", got)
	}
}

// A window Fyne cannot create raises one dialog and ends the process.
func TestFatalWindowErrorRaisesOneDialogAndExits(t *testing.T) {
	var (
		calls  atomic.Int32
		exited atomic.Int32
	)

	previousNotify := notifyStartupError
	previousExit := exitWindow

	notifyStartupError = func(string) { calls.Add(1) }
	exitWindow = func(int) { exited.Add(1) }

	t.Cleanup(func() {
		notifyStartupError = previousNotify
		exitWindow = previousExit
		log.SetOutput(os.Stderr)
	})

	entry := startWindowLog(t.TempDir())

	log.Printf("Fyne error:  window creation error")
	log.Printf("  Cause: APIUnavailable: WGL: The driver does not appear to support OpenGL")

	entry.keep()

	if !waitFor(calls.Load, 1) {
		t.Fatalf("notifications = %d, want one", calls.Load())
	}

	if !waitFor(exited.Load, 1) {
		t.Fatalf("exits = %d, want the windowless process to end", exited.Load())
	}

	time.Sleep(300 * time.Millisecond)

	if got := calls.Load(); got != 1 {
		t.Errorf("notifications = %d after a second line, want one", got)
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

	time.Sleep(400 * time.Millisecond)

	if got := calls.Load(); got != 0 {
		t.Errorf("notifications = %d, want none", got)
	}
}

// A window that cannot be created is fatal; a recoverable Fyne error is not.
func TestIsFatalWindowError(t *testing.T) {
	tests := []struct {
		line  string
		fatal bool
	}{
		{"Fyne error:  window creation error", true},
		{"Fyne error: failed to initialise GLFW", true},
		{"Fyne error: something recoverable", false},
	}

	for _, tt := range tests {
		if got := isFatalWindowError(tt.line); got != tt.fatal {
			t.Errorf("isFatalWindowError(%q) = %v, want %v", tt.line, got, tt.fatal)
		}
	}
}

// waitFor polls until get reports want, so an asynchronous notification does
// not need a fixed sleep in the assertion.
func waitFor(get func() int32, want int32) bool {
	deadline := time.Now().Add(3 * time.Second)

	for time.Now().Before(deadline) {
		if get() == want {
			return true
		}

		time.Sleep(20 * time.Millisecond)
	}

	return false
}
