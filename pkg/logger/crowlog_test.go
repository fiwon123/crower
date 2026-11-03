package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestLogger(t *testing.T) (*Logger, *strings.Builder, *strings.Builder, string) {
	t.Helper()

	logPath := filepath.Join(t.TempDir(), "crower.log")

	out := &strings.Builder{}
	errOut := &strings.Builder{}

	logger, err := New(logPath, out, errOut)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Cleanup(func() { logger.Close() })

	return logger, out, errOut, logPath
}

// The logger takes plain values, not key/value pairs, so every value it is
// given has to reach the terminal. An earlier version walked the values as if
// they were pairs and printed every other one, so a message lost half of what
// it said: "need to pass arguments" arrived without the name of the argument.
func TestLoggerPrintsEveryValue(t *testing.T) {
	logger, _, errOut, _ := newTestLogger(t)

	logger.Info("need to pass arguments:", "path")

	if !strings.Contains(errOut.String(), "path") {
		t.Errorf("expected the value to be printed, got %q", errOut.String())
	}

	if !strings.Contains(errOut.String(), "need to pass arguments:") {
		t.Errorf("expected the message to be printed, got %q", errOut.String())
	}
}

// A result is data a script may pipe: it keeps its exact bytes, so no trailing
// space is added to it. Every result line used to end with one.
func TestOutKeepsTheBytesItIsGiven(t *testing.T) {
	logger, out, _, _ := newTestLogger(t)

	logger.Out("name  exec")

	if got, want := out.String(), "name  exec\n"; got != want {
		t.Errorf("Out() wrote %q, want %q", got, want)
	}
}

// Results and diagnostics must not share a stream: a script that pipes the
// results would otherwise receive the prompts as well.
func TestResultsAndDiagnosticsUseDifferentStreams(t *testing.T) {
	logger, out, errOut, _ := newTestLogger(t)

	logger.Out("a result")
	logger.Info("a diagnostic")

	if strings.Contains(out.String(), "a diagnostic") {
		t.Errorf("diagnostic reached the result stream: %q", out.String())
	}

	if strings.Contains(errOut.String(), "a result") {
		t.Errorf("result reached the diagnostic stream: %q", errOut.String())
	}
}

// Debug is the log's business, not the terminal's.
func TestDebugReachesTheFileOnly(t *testing.T) {
	logger, out, errOut, logPath := newTestLogger(t)

	logger.Debug("sh -c echo hello")

	if out.Len() != 0 {
		t.Errorf("Debug() wrote to the result stream: %q", out.String())
	}

	if errOut.Len() != 0 {
		t.Errorf("Debug() wrote to the diagnostic stream: %q", errOut.String())
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("error reading the log file: %v", err)
	}

	if !strings.Contains(string(content), "sh -c echo hello") {
		t.Errorf("expected the line in the log, got %q", string(content))
	}
}

// The file log gets the values too, and it gets them as a value: slog writes
// "!BADKEY" for a value that is not a key/value pair.
func TestLoggerWritesEveryValueToTheFile(t *testing.T) {
	logger, _, _, logPath := newTestLogger(t)

	logger.Info("need to pass arguments:", "path")

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("error reading the log file: %v", err)
	}

	if !strings.Contains(string(content), "path") {
		t.Errorf("expected the value to be logged, got %q", string(content))
	}

	if strings.Contains(string(content), "!BADKEY") {
		t.Errorf("expected no bad key in the log, got %q", string(content))
	}
}

// A log path that cannot be opened is an error, not a panic: the log is a
// convenience, and a read-only home directory must not turn `crower list` into
// a crash.
func TestNewReportsAnUnwritableLogPath(t *testing.T) {
	_, err := New(filepath.Join(t.TempDir(), "missing", "dir", "crower.log"), nil, nil)
	if err == nil {
		t.Fatal("New() with an unwritable path returned no error")
	}
}
