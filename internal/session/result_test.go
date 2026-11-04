package session_test

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage/command"
	"github.com/fiwon123/crower/internal/session"
)

// newConfig builds a run against buffers instead of a terminal, rooted in a
// fresh temporary directory so a test never writes into the working tree.
func newConfig(t *testing.T, format session.Format) (*session.Session, *bytes.Buffer) {
	t.Helper()

	out := &bytes.Buffer{}

	cfg, err := session.New(
		filepath.Join(t.TempDir(), "crower.yaml"),
		nil,
		command.NewMapData(),
		command.NewMapData(),
		session.Options{
			Streams: session.Streams{In: strings.NewReader(""), Out: out, Err: &bytes.Buffer{}},
			Format:  format,
		},
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	return cfg, out
}

// Text mode must not touch the bytes an operation built: the table alignment is
// the operation's, and rewriting it here would change what a script reads.
func TestResultTextIsVerbatim(t *testing.T) {
	cfg, out := newConfig(t, session.FormatText)

	if err := cfg.Result("row   name\n", map[string]string{"ignored": "yes"}); err != nil {
		t.Fatalf("Result: %v", err)
	}

	if got, want := out.String(), "row   name\n"; got != want {
		t.Errorf("text result = %q, want %q", got, want)
	}
}

func TestResultJSONIsOneObject(t *testing.T) {
	cfg, out := newConfig(t, session.FormatJSON)

	payload := struct {
		Commands []string `json:"commands"`
	}{Commands: []string{"a", "b"}}

	if err := cfg.Result("a table nobody should see", payload); err != nil {
		t.Fatalf("Result: %v", err)
	}

	got := out.String()

	if strings.Count(got, "\n") != 1 || !strings.HasSuffix(got, "\n") {
		t.Fatalf("json result is not exactly one line: %q", got)
	}

	if strings.Contains(got, "nobody") {
		t.Errorf("text rendering leaked into the json result: %q", got)
	}

	var decoded struct {
		Commands []string `json:"commands"`
	}
	if err := json.Unmarshal([]byte(got), &decoded); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, got)
	}

	if len(decoded.Commands) != 2 {
		t.Errorf("commands = %v, want 2 entries", decoded.Commands)
	}
}

// A path is data, not HTML: \u0026 would make the value differ from the path
// the caller passed in.
func TestResultJSONDoesNotEscapeHTML(t *testing.T) {
	cfg, out := newConfig(t, session.FormatJSON)

	payload := struct {
		Path string `json:"path"`
	}{Path: "a&b<c>"}

	if err := cfg.Result("", payload); err != nil {
		t.Fatalf("Result: %v", err)
	}

	if got := out.String(); !strings.Contains(got, "a&b<c>") {
		t.Errorf("path was escaped: %q", got)
	}
}

// An empty collection has to encode as [] so a script does not have to special
// case null.
func TestResultJSONEmptyCollection(t *testing.T) {
	cfg, out := newConfig(t, session.FormatJSON)

	payload := struct {
		Items []string `json:"items"`
	}{Items: []string{}}

	if err := cfg.Result("", payload); err != nil {
		t.Fatalf("Result: %v", err)
	}

	if got := out.String(); !strings.Contains(got, `"items":[]`) {
		t.Errorf("empty collection = %q, want []", got)
	}
}

func TestResultUnknownFormatIsAnError(t *testing.T) {
	cfg, out := newConfig(t, session.Format("yaml"))

	if err := cfg.Result("table", map[string]string{"a": "b"}); err == nil {
		t.Fatal("Result accepted a format it cannot render")
	}

	if out.String() != "" {
		t.Errorf("unknown format wrote a result anyway: %q", out.String())
	}
}

func TestResultJSONEncodeFailureIsAnError(t *testing.T) {
	cfg, _ := newConfig(t, session.FormatJSON)

	if err := cfg.Result("", make(chan int)); err == nil {
		t.Fatal("Result accepted a payload it cannot encode")
	}
}
