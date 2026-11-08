//go:build fyne

package fyneui

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// The Browse button sits a gutter away from its field and is centred on the row,
// instead of wherever the row's height happens to put it (#510).
func TestBrowseButtonSitsBesideAndCentredOnItsField(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	entry := newScrollFreeEntry()
	field := fieldWithBrowse(entry, func() {})

	window := test.NewWindow(container.NewVBox(field))
	t.Cleanup(window.Close)
	window.Resize(fyne.NewSize(600, 120))

	buttons := buttonsIn(field)
	if len(buttons) != 1 {
		t.Fatalf("found %d buttons, want the Browse control", len(buttons))
	}

	button := buttons[0]
	driver := fyne.CurrentApp().Driver()

	entryPos := driver.AbsolutePositionForObject(entry)
	buttonPos := driver.AbsolutePositionForObject(button.Button)

	// The Border pads its own right slot, so the visible gap is that padding plus
	// the gutter the layout was given.
	if gap, want := buttonPos.X-(entryPos.X+entry.Size().Width), theme.Padding()+actionGap(); gap != want {
		t.Errorf("gutter = %.1f, want %v", gap, want)
	}

	entryCentre := entryPos.Y + entry.Size().Height/2
	buttonCentre := buttonPos.Y + button.Button.Size().Height/2

	if diff := buttonCentre - entryCentre; diff < -1 || diff > 1 {
		t.Errorf("button centre = %.1f, want the field centre %.1f", buttonCentre, entryCentre)
	}
}

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

func (s *stubTools) Copy(_ context.Context, _ func(desktop.Progress), sources []string, destination string) (string, timeline.Log, error) {
	s.copies = append(s.copies, [2]any{sources, destination})

	return s.result()
}

func (s *stubTools) Move(_ context.Context, _ func(desktop.Progress), sources []string, destination string) (string, timeline.Log, error) {
	s.moves = append(s.moves, [2]any{sources, destination})

	return s.result()
}

func (s *stubTools) Search(mode desktop.SearchMode, term, path string) (string, timeline.Log, error) {
	s.searches = append(s.searches, [3]any{mode, term, path})

	return s.result()
}

func (s *stubTools) Extract(_ context.Context, _ func(desktop.Progress), archives []string, outDir string) (string, timeline.Log, error) {
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

// The tools tab builds its four cards, each with its own result area, and no
// result is visible before an action ran.
func TestToolsTabBuilds(t *testing.T) {
	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})

	if content := front.toolsTab(); content == nil {
		t.Fatal("toolsTab returned nil")
	}

	results := map[string]*toolResult{
		"transfer": front.transferResult,
		"search":   front.searchResult,
		"extract":  front.extractResult,
		"open":     front.openResult,
	}

	for name, result := range results {
		if result == nil {
			t.Fatalf("%s result was not built", name)
		}

		if result.area.Visible() {
			t.Errorf("%s result is visible before an action ran", name)
		}
	}
}

// Copy and move need a source and a destination before anything runs.
func TestTransferNeedsPathsAndADestination(t *testing.T) {
	tools := &stubTools{}
	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})
	front.tools = tools

	if _, _, err := front.transfer(context.Background(), nil, true, nil, "/tmp"); err == nil {
		t.Error("transfer accepted no source")
	}

	if _, _, err := front.transfer(context.Background(), nil, true, []string{"/tmp/a"}, "  "); err == nil {
		t.Error("transfer accepted no destination")
	}

	if len(tools.copies) != 0 {
		t.Fatalf("copies = %v, want none before validation", tools.copies)
	}

	if _, _, err := front.transfer(context.Background(), nil, true, []string{"/tmp/a"}, "/tmp/out"); err != nil {
		t.Fatalf("transfer: %v", err)
	}

	if len(tools.copies) != 1 {
		t.Fatalf("copies = %v, want one", tools.copies)
	}

	if _, _, err := front.transfer(context.Background(), nil, false, []string{"/tmp/a"}, "/tmp/out"); err != nil {
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

	if _, _, err := front.extract(context.Background(), nil, nil, ""); err == nil {
		t.Error("extract accepted no archive")
	}

	if _, _, err := front.open(nil); err == nil {
		t.Error("open accepted no path")
	}

	if _, _, err := front.extract(context.Background(), nil, []string{"a.zip"}, "/tmp/out"); err != nil {
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

	front.reportTool(front.transferResult, "Copied 1 item(s)", history, nil)

	if len(front.model.History.AllData) != 1 {
		t.Fatalf("model history = %+v, want the recorded log", front.model.History.AllData)
	}

	if front.transferResult.label.Text != "Copied 1 item(s)" {
		t.Errorf("result = %q, want the message", front.transferResult.label.Text)
	}

	if !front.transferResult.area.Visible() || !front.transferResult.body.Visible() {
		t.Error("the result area did not open")
	}

	if front.searchResult.area.Visible() {
		t.Error("a transfer result appeared in the search card")
	}
}

func TestReportToolShowsTheErrorAndKeepsTheHistory(t *testing.T) {
	front, _ := newTestUI(t, testModel(), stubEditor{}, stubReverter{})
	front.toolsTab()

	var history timeline.Log
	history.Add(operation.Copy, "", "note")

	front.reportTool(front.transferResult, "", history, errors.New("permission denied"))

	if !strings.Contains(front.transferResult.label.Text, "permission denied") {
		t.Errorf("result = %q, want the error named", front.transferResult.label.Text)
	}

	if len(front.model.History.AllData) != 0 {
		t.Errorf("history = %+v, want it untouched after a failure", front.model.History.AllData)
	}
}

// The result toggle collapses one card's square without losing its text.
func TestToolResultToggle(t *testing.T) {
	result := newToolResult()

	if result.area.Visible() {
		t.Fatal("a fresh result is visible")
	}

	result.report("done")

	if !result.body.Visible() || result.toggle.Text != "Hide result" {
		t.Error("report did not open the result")
	}

	result.setVisible(false)

	if result.body.Visible() || result.toggle.Text != "Show result" {
		t.Error("the toggle did not collapse the result")
	}

	if result.label.Text != "done" {
		t.Errorf("text = %q, want it kept while collapsed", result.label.Text)
	}
}

// appendPath builds a comma-separated list one Browse pick at a time.
func TestAppendPath(t *testing.T) {
	entry := newScrollFreeEntry()

	appendPath(entry, "/tmp/a")
	appendPath(entry, "/tmp/b")

	if entry.Text != "/tmp/a, /tmp/b" {
		t.Errorf("text = %q, want both paths", entry.Text)
	}

	if got := splitPaths(entry.Text); len(got) != 2 || got[0] != "/tmp/a" || got[1] != "/tmp/b" {
		t.Errorf("splitPaths(%q) = %v, want both paths", entry.Text, got)
	}
}

// A Fyne URI reports a Windows path with a leading slash; the CLI needs the
// path as the file system spells it.
func TestURIPathDropsTheWindowsLeadingSlash(t *testing.T) {
	uri, err := storage.ParseURI("file:///C:/Users/dev/file.txt")
	if err != nil {
		t.Fatalf("ParseURI: %v", err)
	}

	if got, want := uriPath(uri), filepath.FromSlash("C:/Users/dev/file.txt"); got != want {
		t.Errorf("uriPath = %q, want %q", got, want)
	}

	unix, err := storage.ParseURI("file:///tmp/file.txt")
	if err != nil {
		t.Fatalf("ParseURI: %v", err)
	}

	if got, want := uriPath(unix), filepath.FromSlash("/tmp/file.txt"); got != want {
		t.Errorf("uriPath = %q, want %q", got, want)
	}
}

// A cancellable run offers Cancel, which stops the context; the button goes
// away when the run ends.
func TestToolResultCancel(t *testing.T) {
	result := newToolResult()

	cancelled := false
	result.start("Copying…", func() { cancelled = true })

	if !result.cancelButton.Visible() {
		t.Fatal("a cancellable run did not offer Cancel")
	}

	result.cancelRun()

	if !cancelled {
		t.Error("Cancel did not stop the run")
	}

	if !result.cancelButton.Disabled() {
		t.Error("Cancel stayed enabled after it was pressed")
	}

	if result.label.Text != "Cancelling…" {
		t.Errorf("label = %q, want the cancelling note", result.label.Text)
	}

	result.stop()

	if result.cancelButton.Visible() {
		t.Error("Cancel stayed visible after the run ended")
	}
}

// A run that cannot be stopped — a search, an open — offers no Cancel button.
func TestToolResultWithoutCancel(t *testing.T) {
	result := newToolResult()
	result.start("Searching…", nil)

	if result.cancelButton.Visible() {
		t.Error("a non-cancellable run offered Cancel")
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
