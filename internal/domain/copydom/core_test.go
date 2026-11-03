package copydom_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/copydom"
	"github.com/fiwon123/crower/internal/helper"
)

// The history entry has to describe the whole command.
//
// The note was built from the argument list after the destination had been
// taken off the end, so `crower copy a.txt out/` was recorded as
// `crower copy a.txt`: the entry could not say where the file went.
func TestCopyRecordsTheWholeCommand(t *testing.T) {
	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	sourceDir := t.TempDir()
	source := filepath.Join(sourceDir, "a.txt")
	if err := os.WriteFile(source, []byte("hello"), 0o600); err != nil {
		t.Fatalf("error writing the source file: %v", err)
	}

	destination := t.TempDir()

	if err := copydom.NewCore(app).Copy([]string{source, destination}); err != nil {
		t.Fatalf("Copy() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(destination, "a.txt")); err != nil {
		t.Fatalf("the file was not copied: %v", err)
	}

	if len(app.History.AllData) == 0 {
		t.Fatal("the copy recorded no history entry")
	}

	note := app.History.AllData[len(app.History.AllData)-1].Note

	if !strings.Contains(note, "crower copy") {
		t.Errorf("note = %q, want it to name the command", note)
	}

	if strings.Contains(note, "crower crower") {
		t.Errorf("note = %q, want the program name once", note)
	}

	if !strings.Contains(note, destination) {
		t.Errorf("note = %q, want it to name the destination %q", note, destination)
	}

	if strings.HasSuffix(note, " ") {
		t.Errorf("note = %q, want no trailing space", note)
	}
}

// A batch is attempted whole and reports the failures together, so a script can
// see that some of the work did not happen.
func TestCopyReportsEveryFailure(t *testing.T) {
	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	destination := t.TempDir()

	err = copydom.NewCore(app).Copy([]string{
		filepath.Join(t.TempDir(), "missing-1"),
		filepath.Join(t.TempDir(), "missing-2"),
		destination,
	})
	if err == nil {
		t.Fatal("Copy() of two missing paths returned no error")
	}

	message := err.Error()
	if !strings.Contains(message, "missing-1") || !strings.Contains(message, "missing-2") {
		t.Errorf("error = %q, want both failures", message)
	}
}
