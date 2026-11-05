//go:build fyne

package fyneui

import (
	"errors"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// stubTools records what the tools tab asked for.
type stubTools struct {
	copies   [][2]any
	moves    [][2]any
	searches [][3]any
	extracts [][2]any
	opens    [][]string
	history  timeline.Log
	err      error
}

func (s *stubTools) Copy(sources []string, destination string) (string, timeline.Log, error) {
	s.copies = append(s.copies, [2]any{sources, destination})

	return s.result()
}

func (s *stubTools) Move(sources []string, destination string) (string, timeline.Log, error) {
	s.moves = append(s.moves, [2]any{sources, destination})

	return s.result()
}

func (s *stubTools) Search(mode desktop.SearchMode, term, path string) (string, timeline.Log, error) {
	s.searches = append(s.searches, [3]any{mode, term, path})

	return s.result()
}

func (s *stubTools) Extract(archives []string, outDir string) (string, timeline.Log, error) {
	s.extracts = append(s.extracts, [2]any{archives, outDir})

	return s.result()
}

func (s *stubTools) Open(paths []string) (string, timeline.Log, error) {
	s.opens = append(s.opens, paths)

	return s.result()
}

func (s *stubTools) result() (string, timeline.Log, error) {
	if s.err != nil {
		return "", timeline.Log{}, s.err
	}

	return "done", s.history, nil
}

// The tools tab builds its four cards and the shared result pane.
func TestToolsTabBuilds(t *testing.T) {
	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})

	if content := front.toolsTab(); content == nil {
		t.Fatal("toolsTab returned nil")
	}

	if front.toolResult == nil {
		t.Fatal("the result pane was not built")
	}

	if !strings.Contains(front.toolResult.Text, "Pick a tool") {
		t.Errorf("result = %q, want the placeholder", front.toolResult.Text)
	}
}

// Copy and move need a source and a destination before anything runs.
func TestTransferNeedsPathsAndADestination(t *testing.T) {
	tools := &stubTools{}
	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})
	front.tools = tools

	if _, _, err := front.transfer(true, nil, "/tmp"); err == nil {
		t.Error("transfer accepted no source")
	}

	if _, _, err := front.transfer(true, []string{"/tmp/a"}, "  "); err == nil {
		t.Error("transfer accepted no destination")
	}

	if len(tools.copies) != 0 {
		t.Fatalf("copies = %v, want none before validation", tools.copies)
	}

	if _, _, err := front.transfer(true, []string{"/tmp/a"}, "/tmp/out"); err != nil {
		t.Fatalf("transfer: %v", err)
	}

	if len(tools.copies) != 1 {
		t.Fatalf("copies = %v, want one", tools.copies)
	}

	if _, _, err := front.transfer(false, []string{"/tmp/a"}, "/tmp/out"); err != nil {
		t.Fatalf("move: %v", err)
	}

	if len(tools.moves) != 1 {
		t.Fatalf("moves = %v, want one", tools.moves)
	}
}

// A search needs a term; the folder defaults to the current one, and a web
// search never carries a path.
func TestSearchDefaultsThePath(t *testing.T) {
	tools := &stubTools{}
	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})
	front.tools = tools

	if _, _, err := front.search(desktop.SearchBoth, "  ", "."); err == nil {
		t.Error("search accepted an empty term")
	}

	if _, _, err := front.search(desktop.SearchFiles, "needle", ""); err != nil {
		t.Fatalf("search: %v", err)
	}

	first := tools.searches[0]
	if first[0] != desktop.SearchFiles || first[1] != "needle" || first[2] != "." {
		t.Errorf("search call = %v, want files/needle/.", first)
	}

	if _, _, err := front.search(desktop.SearchWeb, "golang", "/tmp"); err != nil {
		t.Fatalf("web search: %v", err)
	}

	second := tools.searches[1]
	if second[2] != "" {
		t.Errorf("web search path = %v, want it empty", second[2])
	}
}

func TestExtractAndOpenNeedPaths(t *testing.T) {
	tools := &stubTools{}
	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})
	front.tools = tools

	if _, _, err := front.extract(nil, ""); err == nil {
		t.Error("extract accepted no archive")
	}

	if _, _, err := front.open(nil); err == nil {
		t.Error("open accepted no path")
	}

	if _, _, err := front.extract([]string{"a.zip"}, "/tmp/out"); err != nil {
		t.Fatalf("extract: %v", err)
	}

	if _, _, err := front.open([]string{"/tmp/a"}); err != nil {
		t.Fatalf("open: %v", err)
	}

	if len(tools.extracts) != 1 || len(tools.opens) != 1 {
		t.Fatalf("extracts/opens = %v/%v, want one each", tools.extracts, tools.opens)
	}
}

// A successful tool run puts the history it recorded into the model, so the
// History tab lists it without reopening the window.
func TestReportToolReplacesTheHistory(t *testing.T) {
	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})
	front.toolsTab()

	var history timeline.Log
	history.Add(operation.Copy, "", "note")

	front.reportTool("Copied 1 item(s)", history, nil)

	if len(front.model.History.AllData) != 1 {
		t.Fatalf("model history = %+v, want the recorded log", front.model.History.AllData)
	}

	if front.toolResult.Text != "Copied 1 item(s)" {
		t.Errorf("result = %q, want the message", front.toolResult.Text)
	}
}

func TestReportToolShowsTheErrorAndKeepsTheHistory(t *testing.T) {
	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})
	front.toolsTab()

	var history timeline.Log
	history.Add(operation.Copy, "", "note")

	front.reportTool("", history, errors.New("permission denied"))

	if !strings.Contains(front.toolResult.Text, "permission denied") {
		t.Errorf("result = %q, want the error named", front.toolResult.Text)
	}

	if len(front.model.History.AllData) != 0 {
		t.Errorf("history = %+v, want it untouched after a failure", front.model.History.AllData)
	}
}

// splitPaths accepts a pasted shell line and one path per line.
func TestSplitPaths(t *testing.T) {
	got := splitPaths(" /tmp/a \n/tmp/b\n\n/tmp/c,/tmp/d ")

	want := []string{"/tmp/a", "/tmp/b", "/tmp/c", "/tmp/d"}
	if len(got) != len(want) {
		t.Fatalf("splitPaths = %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("path %d = %q, want %q", i, got[i], want[i])
		}
	}
}
