package extractcontent_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/extractcontent"
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

type fakeOutput struct {
	results []string
}

func (o *fakeOutput) Result(text string) {
	o.results = append(o.results, text)
}

type fakeArchives struct {
	calls   []string
	out     string
	err     error
	errOn   string
	outDirs []string
}

func (a *fakeArchives) Extract(ext, archive, outDir string) (string, error) {
	a.calls = append(a.calls, ext+" "+archive+" "+outDir)
	a.outDirs = append(a.outDirs, outDir)

	if a.errOn != "" && archive == a.errOn {
		return a.out, a.err
	}

	return a.out, nil
}

type fakeNotes struct{}

func (fakeNotes) ExtractNote(out string, args []string) string {
	return "note:" + out + ":" + strings.Join(args, ",")
}

type harness struct {
	state    *fakeState
	output   *fakeOutput
	archives *fakeArchives
}

func newHarness() *harness {
	return &harness{state: &fakeState{}, output: &fakeOutput{}, archives: &fakeArchives{}}
}

func (h *harness) run(request extractcontent.Request) error {
	return extractcontent.New(h.state, h.output, h.archives, fakeNotes{}).Run(request)
}

// A cancelled context stops the batch before the first archive, and progress
// names each archive as it starts.
func TestExtractReportsProgressAndStopsOnCancel(t *testing.T) {
	reported := make([]string, 0)

	h := newHarness()
	err := h.run(extractcontent.Request{
		Args: []string{"a.zip", "b.tar"},
		Progress: func(done, total int, current string) {
			reported = append(reported, current)
		},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !slices.Equal(reported, []string{"a.zip", "b.tar"}) {
		t.Errorf("progress = %v, want each archive", reported)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	h = newHarness()

	err = h.run(extractcontent.Request{Args: []string{"a.zip"}, Context: ctx})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}

	if len(h.archives.calls) != 0 {
		t.Errorf("calls = %v, want none after a cancel", h.archives.calls)
	}
}

func TestExtractGlobsAndUsesEachArchivesFolder(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "a.zip")
	second := filepath.Join(dir, "b.tar")

	for _, path := range []string{first, second} {
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	h := newHarness()

	if err := h.run(extractcontent.Request{Args: []string{filepath.Join(dir, "*")}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.archives.calls) != 2 {
		t.Fatalf("calls = %v, want both archives", h.archives.calls)
	}

	for _, outDir := range h.archives.outDirs {
		if outDir != dir {
			t.Errorf("outDir = %q, want the archive's own folder %q", outDir, dir)
		}
	}

	if h.state.records[0].op != operation.Extract {
		t.Errorf("recorded %q, want extract", h.state.records[0].op)
	}

	if h.state.records[0].note != "note::"+filepath.Join(dir, "*") {
		t.Errorf("note = %q, want the original pattern", h.state.records[0].note)
	}
}

// Each archive keeps its own default destination: two archives in different
// folders must not both go to the first one's.
func TestExtractKeepsAPerArchiveDestination(t *testing.T) {
	one := t.TempDir()
	two := t.TempDir()

	first := filepath.Join(one, "a.zip")
	second := filepath.Join(two, "b.zip")

	for _, path := range []string{first, second} {
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	h := newHarness()

	if err := h.run(extractcontent.Request{Args: []string{first, second}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.archives.outDirs, []string{one, two}; !slices.Equal(got, want) {
		t.Errorf("outDirs = %v, want %v", got, want)
	}
}

func TestExtractHonoursOutDir(t *testing.T) {
	h := newHarness()

	if err := h.run(extractcontent.Request{Args: []string{"a.zip", "b.tar"}, OutDir: "./build"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	for _, outDir := range h.archives.outDirs {
		if outDir != "./build" {
			t.Errorf("outDir = %q, want ./build", outDir)
		}
	}

	if h.state.records[0].note != "note:./build:a.zip,b.tar" {
		t.Errorf("note = %q, want the out dir recorded", h.state.records[0].note)
	}
}

func TestExtractWithoutAnExtension(t *testing.T) {
	h := newHarness()

	err := h.run(extractcontent.Request{Args: []string{"archive."}})
	if err == nil || !strings.Contains(err.Error(), "no file extension") {
		t.Fatalf("error = %v, want the extension error", err)
	}

	if len(h.state.records) != 0 {
		t.Error("a failed extraction recorded an entry")
	}
}

func TestExtractJoinsFailures(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness()
	h.archives.err = boom
	h.archives.errOn = "a.zip"

	err := h.run(extractcontent.Request{Args: []string{"a.zip", "b.zip"}})
	if !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}

	if len(h.archives.calls) != 2 {
		t.Errorf("calls = %v, want every archive attempted", h.archives.calls)
	}

	if len(h.state.records) != 0 {
		t.Error("a partially failed extraction recorded an entry")
	}
}

func TestExtractWritesTheToolOutputTrimmed(t *testing.T) {
	h := newHarness()
	h.archives.out = "expanded\nanother\n"

	if err := h.run(extractcontent.Request{Args: []string{"a.zip"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.output.results, []string{"expanded\nanother"}; !slices.Equal(got, want) {
		t.Errorf("results = %v, want %v", got, want)
	}
}

func TestExtractWithNoPaths(t *testing.T) {
	h := newHarness()

	err := h.run(extractcontent.Request{})
	if err == nil || !strings.Contains(err.Error(), "no paths given") {
		t.Fatalf("error = %v, want the empty-paths error", err)
	}
}

func TestExtractRecordErrorFailsTheRun(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness()
	h.state.recordErr = boom

	if err := h.run(extractcontent.Request{Args: []string{"a.zip"}}); !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}
}
