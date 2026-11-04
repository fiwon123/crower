package opencontent_test

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/opencontent"
	"github.com/fiwon123/crower/internal/domain/operation"
)

type fakeState struct {
	configFile   string
	configFolder string
	records      []record
	recordErr    error
}

type record struct {
	op          operation.Operation
	commandName string
	note        string
}

func (s *fakeState) ConfigFile() string {
	return s.configFile
}

func (s *fakeState) ConfigFolder() string {
	return s.configFolder
}

func (s *fakeState) Record(op operation.Operation, commandName, note string) error {
	if s.recordErr != nil {
		return s.recordErr
	}

	s.records = append(s.records, record{op: op, commandName: commandName, note: note})
	return nil
}

type fakeOutput struct {
	infos   []string
	results []string
}

func (o *fakeOutput) Info(message string, data ...any) {
	o.infos = append(o.infos, message)
}

func (o *fakeOutput) Result(text string) {
	o.results = append(o.results, text)
}

type fakeOpener struct {
	opened  []string
	command string
	output  string
	err     error
	errOn   string
}

func (o *fakeOpener) Open(path string) (string, string, error) {
	o.opened = append(o.opened, path)

	if o.errOn != "" && path == o.errOn {
		return o.command, o.output, o.err
	}

	return o.command, o.output, nil
}

type fakeEditor struct {
	calls int
	err   error
}

func (e *fakeEditor) OpenEditor() error {
	e.calls++
	return e.err
}

type fakeNotes struct{}

func (fakeNotes) OpenNote(args []string) string     { return "note:open:" + strings.Join(args, ",") }
func (fakeNotes) OpenFileNote(args []string) string { return "note:file:" + strings.Join(args, ",") }
func (fakeNotes) OpenFolderNote(args []string) string {
	return "note:folder:" + strings.Join(args, ",")
}
func (fakeNotes) OpenSystemNote() string { return "note:system" }

type harness struct {
	state  *fakeState
	output *fakeOutput
	opener *fakeOpener
	editor *fakeEditor
}

func newHarness() *harness {
	return &harness{
		state:  &fakeState{configFile: "/cfg/crower.yaml", configFolder: "/cfg"},
		output: &fakeOutput{},
		opener: &fakeOpener{command: "xdg-open 'x'"},
		editor: &fakeEditor{},
	}
}

func (h *harness) run(request opencontent.Request) error {
	return opencontent.New(h.state, h.output, h.opener, h.editor, fakeNotes{}).Run(request)
}

func TestOpenPathsExpandsGlobs(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "a.txt")
	second := filepath.Join(dir, "b.txt")

	for _, path := range []string{first, second} {
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	h := newHarness()

	if err := h.run(opencontent.Request{Mode: opencontent.ModePaths, Args: []string{filepath.Join(dir, "*.txt")}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.opener.opened) != 2 {
		t.Fatalf("opened %v, want both matches", h.opener.opened)
	}

	if !slices.Contains(h.opener.opened, first) || !slices.Contains(h.opener.opened, second) {
		t.Errorf("opened %v, want %s and %s", h.opener.opened, first, second)
	}

	if h.state.records[0].note != "note:open:"+filepath.Join(dir, "*.txt") {
		t.Errorf("note = %q, want the original pattern recorded", h.state.records[0].note)
	}
}

// An argument that matches nothing is opened as written, and a broken pattern
// is skipped rather than failing the run.
func TestOpenPathsKeepsUnmatchedArgumentsAndSkipsBadPatterns(t *testing.T) {
	h := newHarness()

	args := []string{"literal.txt", "["}

	if err := h.run(opencontent.Request{Mode: opencontent.ModePaths, Args: args}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	want := []string{"literal.txt"}
	if !slices.Equal(h.opener.opened, want) {
		t.Errorf("opened %v, want %v", h.opener.opened, want)
	}
}

func TestOpenDefaults(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		h := newHarness()

		if err := h.run(opencontent.Request{Mode: opencontent.ModeFile}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if got, want := h.opener.opened, []string{"/cfg/crower.yaml"}; !slices.Equal(got, want) {
			t.Errorf("opened %v, want %v", got, want)
		}

		if h.state.records[0].note != "note:file:/cfg/crower.yaml" {
			t.Errorf("note = %q, want the file note with the default", h.state.records[0].note)
		}
	})

	t.Run("folder", func(t *testing.T) {
		h := newHarness()

		if err := h.run(opencontent.Request{Mode: opencontent.ModeFolder}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if got, want := h.opener.opened, []string{"/cfg"}; !slices.Equal(got, want) {
			t.Errorf("opened %v, want %v", got, want)
		}

		if h.state.records[0].note != "note:folder:/cfg" {
			t.Errorf("note = %q, want the folder note with the default", h.state.records[0].note)
		}
	})

	t.Run("explicit paths skip the defaults", func(t *testing.T) {
		h := newHarness()

		if err := h.run(opencontent.Request{Mode: opencontent.ModeFile, Args: []string{"notes.txt"}}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if got, want := h.opener.opened, []string{"notes.txt"}; !slices.Equal(got, want) {
			t.Errorf("opened %v, want %v", got, want)
		}
	})
}

func TestOpenReportsTheCommandAndOutput(t *testing.T) {
	h := newHarness()
	h.opener.output = "opened in browser"

	if err := h.run(opencontent.Request{Mode: opencontent.ModePaths, Args: []string{"a.txt"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.output.infos, []string{"opening: "}; !slices.Equal(got, want) {
		t.Errorf("diagnostics = %v, want the opening line", got)
	}

	if got, want := h.output.results, []string{"opened in browser"}; !slices.Equal(got, want) {
		t.Errorf("results = %v, want the output", got)
	}
}

func TestOpenSilentOutput(t *testing.T) {
	h := newHarness()
	h.opener.output = ""

	if err := h.run(opencontent.Request{Mode: opencontent.ModePaths, Args: []string{"a.txt"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.output.results) != 0 {
		t.Errorf("results = %v, want nothing for empty output", h.output.results)
	}
}

func TestOpenStopsAtTheFirstFailure(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness()
	h.opener.err = boom
	h.opener.errOn = "b.txt"

	err := h.run(opencontent.Request{Mode: opencontent.ModePaths, Args: []string{"a.txt", "b.txt", "c.txt"}})
	if !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}

	if got, want := h.opener.opened, []string{"a.txt", "b.txt"}; !slices.Equal(got, want) {
		t.Errorf("opened %v, want the batch stopped after the failure", got)
	}

	if len(h.state.records) != 0 {
		t.Error("a failed open recorded an entry")
	}
}

func TestOpenSystem(t *testing.T) {
	h := newHarness()

	if err := h.run(opencontent.Request{Mode: opencontent.ModeSystem}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.editor.calls != 1 {
		t.Errorf("editor calls = %d, want 1", h.editor.calls)
	}

	if len(h.opener.opened) != 0 {
		t.Errorf("opened paths = %v, want none", h.opener.opened)
	}

	if h.state.records[0].note != "note:system" {
		t.Errorf("note = %q, want the system note", h.state.records[0].note)
	}
}

func TestOpenErrorPaths(t *testing.T) {
	boom := errors.New("boom")

	t.Run("editor", func(t *testing.T) {
		h := newHarness()
		h.editor.err = boom

		if err := h.run(opencontent.Request{Mode: opencontent.ModeSystem}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.state.records) != 0 {
			t.Error("a failed editor open recorded an entry")
		}
	})

	t.Run("record", func(t *testing.T) {
		h := newHarness()
		h.state.recordErr = boom

		if err := h.run(opencontent.Request{Mode: opencontent.ModePaths, Args: []string{"a.txt"}}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("unknown mode", func(t *testing.T) {
		h := newHarness()

		err := h.run(opencontent.Request{Mode: opencontent.Mode(99)})
		if err == nil || !strings.Contains(err.Error(), "unknown open mode") {
			t.Fatalf("error = %v, want an unknown-mode error", err)
		}
	})
}

func TestOpenRecordsTheOperation(t *testing.T) {
	h := newHarness()

	if err := h.run(opencontent.Request{Mode: opencontent.ModePaths, Args: []string{"a.txt"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	got := h.state.records[0]
	if got.op != operation.Open || got.commandName != "" {
		t.Errorf("recorded %+v, want open with no command name", got)
	}
}
