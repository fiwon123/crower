package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The logger takes plain values, not key/value pairs, so every value it is
// given has to reach the terminal. The helper walked the values as if they
// were pairs and printed every other one, so a message lost half of what it
// said: "need to pass arguments" arrived without the name of the argument.
func TestLoggerPrintsEveryValue(t *testing.T) {
	dir := t.TempDir()

	logger := New(filepath.Join(dir, "crower.log"))

	output := captureStdout(t, func() {
		logger.Info("need to pass arguments:", "path")
	})

	if !strings.Contains(output, "path") {
		t.Errorf("expected the value to be printed, got %q", output)
	}

	if !strings.Contains(output, "need to pass arguments") {
		t.Errorf("expected the message to be printed, got %q", output)
	}
}

// The file log gets the values too, and it gets them as a key/value pair:
// slog writes "!BADKEY" for a value that is not one.
func TestLoggerWritesEveryValueToTheFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "crower.log")

	logger := New(logPath)

	captureStdout(t, func() {
		logger.Info("need to pass arguments:", "path")
	})

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

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("error creating the pipe: %v", err)
	}

	os.Stdout = writer

	fn()

	writer.Close()
	os.Stdout = original

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("error reading the output: %v", err)
	}

	return string(content)
}
