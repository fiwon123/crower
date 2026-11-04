package searchcontent_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/searchcontent"
	"github.com/fiwon123/crower/internal/domain/operation"
)

type fakeState struct {
	records   []record
	recordErr error
}

type record struct {
	op          operation.Operation
	commandName string
	note        string
}

func (s *fakeState) Record(op operation.Operation, commandName, note string) error {
	if s.recordErr != nil {
		return s.recordErr
	}

	s.records = append(s.records, record{op: op, commandName: commandName, note: note})
	return nil
}

type emitted struct {
	text    string
	payload any
}

type fakeOutput struct {
	emitted []emitted
	err     error
}

func (o *fakeOutput) Emit(text string, payload any) error {
	if o.err != nil {
		return o.err
	}

	o.emitted = append(o.emitted, emitted{text: text, payload: payload})
	return nil
}

type fakeSearcher struct {
	calls   []string
	out     string
	err     error
	browser []string
}

func (s *fakeSearcher) Browser(query string) error {
	s.browser = append(s.browser, query)
	return s.err
}

func (s *fakeSearcher) Files(path, term string) (string, error) {
	s.calls = append(s.calls, "files "+path+" "+term)
	return s.out, s.err
}

func (s *fakeSearcher) Folders(path, term string) (string, error) {
	s.calls = append(s.calls, "folders "+path+" "+term)
	return s.out, s.err
}

func (s *fakeSearcher) Both(path, term string) (string, error) {
	s.calls = append(s.calls, "both "+path+" "+term)
	return s.out, s.err
}

type fakeNotes struct{}

func (fakeNotes) SearchBrowserNote([]string) string { return "note:browser" }
func (fakeNotes) SearchFileNote([]string) string    { return "note:file" }
func (fakeNotes) SearchFolderNote([]string) string  { return "note:folder" }
func (fakeNotes) SearchBothNote([]string) string    { return "note:both" }

type harness struct {
	state  *fakeState
	output *fakeOutput
	search *fakeSearcher
}

func newHarness() *harness {
	return &harness{state: &fakeState{}, output: &fakeOutput{}, search: &fakeSearcher{}}
}

func (h *harness) run(request searchcontent.Request) error {
	return searchcontent.New(h.state, h.output, h.search, fakeNotes{}).Run(request)
}

func TestSearchModes(t *testing.T) {
	tests := []struct {
		name      string
		mode      searchcontent.Mode
		args      []string
		wantCall  string
		wantKind  string
		wantQuery string
		wantPath  string
	}{
		{"files", searchcontent.ModeFile, []string{"a.txt", "/tmp"}, "files /tmp a.txt", "file", "a.txt", "/tmp"},
		{"folders", searchcontent.ModeFolder, []string{"logs", "/tmp"}, "folders /tmp logs", "folder", "logs", "/tmp"},
		{"both", searchcontent.ModeBoth, []string{"needle", "/tmp"}, "both /tmp needle", "file_and_folder", "needle", "/tmp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHarness()
			h.search.out = "one\ntwo\n"

			if err := h.run(searchcontent.Request{Mode: tt.mode, Args: tt.args}); err != nil {
				t.Fatalf("Run: %v", err)
			}

			if got, want := h.search.calls, []string{tt.wantCall}; !slices.Equal(got, want) {
				t.Errorf("calls = %v, want %v", got, want)
			}

			if h.output.emitted[0].text != "one\ntwo\n" {
				t.Errorf("text = %q, want the shell output", h.output.emitted[0].text)
			}

			payload := h.output.emitted[0].payload.(searchcontent.Result)
			if payload.Search.Kind != tt.wantKind || payload.Search.Query != tt.wantQuery || payload.Search.Path != tt.wantPath {
				t.Errorf("payload = %+v, want kind/query/path", payload.Search)
			}

			if !slices.Equal(payload.Search.Results, []string{"one", "two"}) {
				t.Errorf("results = %v, want the lines", payload.Search.Results)
			}
		})
	}
}

// The one-positional-value case keeps the legacy parsing: the query stays empty
// and the path stays `.`. It is tracked as a bug; the refactor must not change
// it silently.
func TestSearchWithOneArgumentKeepsTheLegacyParsing(t *testing.T) {
	h := newHarness()

	if err := h.run(searchcontent.Request{Mode: searchcontent.ModeBoth, Args: []string{"needle"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.search.calls, []string{"both . "}; !slices.Equal(got, want) {
		t.Errorf("calls = %v, want %v", got, want)
	}

	payload := h.output.emitted[0].payload.(searchcontent.Result)
	if payload.Search.Query != "" || payload.Search.Path != "." {
		t.Errorf("payload = %+v, want an empty query and the current path", payload.Search)
	}
}

func TestSearchBrowser(t *testing.T) {
	h := newHarness()

	if err := h.run(searchcontent.Request{Mode: searchcontent.ModeBrowser, Args: []string{"golang context"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.search.browser, []string{"golang context"}; !slices.Equal(got, want) {
		t.Errorf("browser = %v, want %v", got, want)
	}

	if h.output.emitted[0].text != "" {
		t.Errorf("text = %q, want empty", h.output.emitted[0].text)
	}

	payload := h.output.emitted[0].payload.(searchcontent.Result)
	if payload.Search.Kind != "browser" || payload.Search.Query != "golang context" || payload.Search.Path != "" {
		t.Errorf("payload = %+v, want a browser result", payload.Search)
	}

	if payload.Search.Results == nil || len(payload.Search.Results) != 0 {
		t.Errorf("results = %#v, want an empty slice", payload.Search.Results)
	}

	if len(h.state.records) != 1 || h.state.records[0].note != "note:browser" {
		t.Errorf("records = %+v, want the browser note", h.state.records)
	}
}

func TestSearchRecordsTheOperation(t *testing.T) {
	h := newHarness()

	if err := h.run(searchcontent.Request{Mode: searchcontent.ModeFile}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	got := h.state.records[0]
	if got.op != operation.Search || got.commandName != "" || got.note != "note:file" {
		t.Errorf("recorded %+v, want search with no command name", got)
	}
}

func TestSearchErrors(t *testing.T) {
	boom := errors.New("boom")

	t.Run("searcher", func(t *testing.T) {
		h := newHarness()
		h.search.err = boom

		if err := h.run(searchcontent.Request{Mode: searchcontent.ModeFile}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.output.emitted) != 0 || len(h.state.records) != 0 {
			t.Error("a failed search still wrote or recorded something")
		}
	})

	t.Run("browser", func(t *testing.T) {
		h := newHarness()
		h.search.err = boom

		if err := h.run(searchcontent.Request{Mode: searchcontent.ModeBrowser}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("emit", func(t *testing.T) {
		h := newHarness()
		h.output.err = boom

		if err := h.run(searchcontent.Request{Mode: searchcontent.ModeBoth}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.state.records) != 0 {
			t.Error("a failed result still recorded an entry")
		}
	})

	t.Run("record", func(t *testing.T) {
		h := newHarness()
		h.state.recordErr = boom

		if err := h.run(searchcontent.Request{Mode: searchcontent.ModeBoth}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}

func TestSearchUnknownMode(t *testing.T) {
	h := newHarness()

	err := h.run(searchcontent.Request{Mode: searchcontent.Mode(99)})
	if err == nil || !strings.Contains(err.Error(), "unknown search mode") {
		t.Fatalf("error = %v, want an unknown-mode error", err)
	}
}

// A CRLF shell output keeps its entries, and an empty line is not a result.
func TestSplitResults(t *testing.T) {
	h := newHarness()
	h.search.out = "a\r\n\r\nb\r\n"

	if err := h.run(searchcontent.Request{Mode: searchcontent.ModeBoth, Args: []string{"x", "."}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	payload := h.output.emitted[0].payload.(searchcontent.Result)
	if !slices.Equal(payload.Search.Results, []string{"a", "b"}) {
		t.Errorf("results = %v, want a and b", payload.Search.Results)
	}
}
