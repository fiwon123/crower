package app

import (
	"bufio"
	"strings"
	"testing"
)

func newReaderConfig(t *testing.T, input string) *Config {
	t.Helper()

	return &Config{in: bufio.NewReader(strings.NewReader(input))}
}

func TestReadLineReturnsTheLine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain", "hello\n", "hello"},
		{"last line without a newline", "hello", "hello"},
		{"empty line", "\n", ""},
		{"windows line ending", "hello\r\n", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newReaderConfig(t, tt.input)

			got, err := cfg.ReadLine()
			if err != nil {
				t.Fatalf("ReadLine() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("ReadLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

// A closed stream is an error, not an empty answer.
//
// The reader used to ignore the error and hand back "", so every prompt in the
// program saw a valid empty answer and asked again: `crower create </dev/null`
// wrote 876,468 lines in five seconds and never exited.
func TestReadLineReportsAClosedStream(t *testing.T) {
	cfg := newReaderConfig(t, "")

	if _, err := cfg.ReadLine(); err == nil {
		t.Fatal("ReadLine() on an empty stream returned no error")
	}
}

// The partial line that arrives before the end of the stream is still a line.
// It is returned once, and the next read reports the end.
func TestReadLineReturnsTheLastLineThenReportsTheEnd(t *testing.T) {
	cfg := newReaderConfig(t, "hello")

	got, err := cfg.ReadLine()
	if err != nil {
		t.Fatalf("ReadLine() error = %v", err)
	}

	if got != "hello" {
		t.Errorf("ReadLine() = %q, want %q", got, "hello")
	}

	if _, err := cfg.ReadLine(); err == nil {
		t.Fatal("the second ReadLine() returned no error after the end of the stream")
	}
}

// Several prompts in one command share one buffer: wrapping the stream again
// for each prompt discarded whatever the previous prompt had already read.
func TestReaderKeepsOneBuffer(t *testing.T) {
	cfg := newReaderConfig(t, "first\nsecond\n")

	if _, err := cfg.ReadLine(); err != nil {
		t.Fatalf("ReadLine() error = %v", err)
	}

	if got, err := cfg.Reader().ReadString('\n'); err != nil || got != "second\n" {
		t.Errorf("Reader() read %q, %v; want %q", got, err, "second\n")
	}
}

// A stream that was left out falls back to the terminal, so a caller that only
// cares about one stream does not have to build the other two.
func TestStreamsFillsInTheDefaults(t *testing.T) {
	filled := Streams{}.withDefaults()

	if filled.In == nil || filled.Out == nil || filled.Err == nil {
		t.Fatalf("withDefaults() left a stream nil: %+v", filled)
	}
}
