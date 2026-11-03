package output_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/output"
)

// A result has to be exactly one line: `crower list --json | wc -l` is 1.
func TestEncodeJSONIsOneLine(t *testing.T) {
	encoded, err := output.EncodeJSON(map[string]any{"commands": []string{"a", "b"}})
	if err != nil {
		t.Fatalf("EncodeJSON: %v", err)
	}

	if strings.HasSuffix(encoded, "\n") {
		t.Errorf("EncodeJSON left a trailing newline: %q", encoded)
	}

	if strings.Contains(encoded, "\n") {
		t.Errorf("EncodeJSON produced more than one line: %q", encoded)
	}
}

// A path or a command is data. Escaping `&` as \u0026 keeps the JSON valid but
// makes the value impossible to compare with the path a script passed in.
func TestEncodeJSONDoesNotEscapeHTML(t *testing.T) {
	encoded, err := output.EncodeJSON(map[string]string{"path": "a&b<c>"})
	if err != nil {
		t.Fatalf("EncodeJSON: %v", err)
	}

	if !strings.Contains(encoded, "a&b<c>") {
		t.Errorf("EncodeJSON escaped HTML characters: %q", encoded)
	}
}

func TestEncodeJSONRejectsAnUnencodableValue(t *testing.T) {
	if _, err := output.EncodeJSON(make(chan int)); err == nil {
		t.Fatal("EncodeJSON of a channel returned no error")
	}
}

// A nil stream is the terminal, so a caller that only cares about one stream
// does not have to construct the other two.
func TestStreamsWithDefaults(t *testing.T) {
	filled := output.Streams{}.WithDefaults()

	if filled.In == nil || filled.Out == nil || filled.Err == nil {
		t.Fatalf("WithDefaults left a stream nil: %+v", filled)
	}
}
