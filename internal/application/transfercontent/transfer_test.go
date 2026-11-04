package transfercontent_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/fiwon123/crower/internal/application/transfercontent"
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

type fakeFiles struct {
	existing map[string]bool
}

func (f *fakeFiles) IsFile(path string) bool {
	return f.existing[path]
}

type fakeTransfer struct {
	calls   []string
	err     error
	errOn   string
	lastDst string
}

func (t *fakeTransfer) Copy(source, destination string, folder bool) error {
	t.lastDst = destination

	kind := "file"
	if folder {
		kind = "folder"
	}

	t.calls = append(t.calls, "copy "+kind+" "+source)

	if t.errOn != "" && source == t.errOn {
		return t.err
	}

	return nil
}

func (t *fakeTransfer) Move(source, destination string) error {
	t.lastDst = destination
	t.calls = append(t.calls, "move "+source)

	if t.errOn != "" && source == t.errOn {
		return t.err
	}

	return nil
}

type fakeNotes struct{}

func (fakeNotes) CopyNote(sub operation.SubOperation, args []string) string {
	return "note:copy:" + string(sub)
}

func (fakeNotes) MoveNote(sub operation.SubOperation, args []string) string {
	return "note:move:" + string(sub)
}

type harness struct {
	state    *fakeState
	files    *fakeFiles
	transfer *fakeTransfer
}

func newHarness(existing ...string) *harness {
	files := map[string]bool{}
	for _, path := range existing {
		files[path] = true
	}

	return &harness{
		state:    &fakeState{},
		files:    &fakeFiles{existing: files},
		transfer: &fakeTransfer{},
	}
}

func (h *harness) run(request transfercontent.Request) error {
	return transfercontent.New(h.state, h.files, h.transfer, fakeNotes{}).Run(request)
}

func TestCopyDetectsFilesAndFolders(t *testing.T) {
	tests := []struct {
		name     string
		existing []string
		args     []string
		want     []string
		wantSub  string
	}{
		{
			"files",
			[]string{"a.txt", "b.txt"},
			[]string{"a.txt", "b.txt", "dest"},
			[]string{"copy file a.txt", "copy file b.txt"},
			"file",
		},
		{
			"folders",
			nil,
			[]string{"one", "two", "dest"},
			[]string{"copy folder one", "copy folder two"},
			"folder",
		},
		{
			"mixed",
			[]string{"a.txt"},
			[]string{"a.txt", "dir", "dest"},
			[]string{"copy file a.txt", "copy folder dir"},
			"file_and_folder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHarness(tt.existing...)

			if err := h.run(transfercontent.Request{Op: operation.Copy, Args: tt.args}); err != nil {
				t.Fatalf("Run: %v", err)
			}

			if !slices.Equal(h.transfer.calls, tt.want) {
				t.Errorf("calls = %v, want %v", h.transfer.calls, tt.want)
			}

			if h.transfer.lastDst != "dest" {
				t.Errorf("destination = %q, want the last argument", h.transfer.lastDst)
			}

			if h.state.records[0].note != "note:copy:"+tt.wantSub {
				t.Errorf("note = %q, want the %s note", h.state.records[0].note, tt.wantSub)
			}

			if h.state.records[0].op != operation.Copy {
				t.Errorf("recorded %q, want copy", h.state.records[0].op)
			}
		})
	}
}

func TestMoveUsesTheMoveTransfer(t *testing.T) {
	h := newHarness("a.txt")

	if err := h.run(transfercontent.Request{Op: operation.Move, Args: []string{"a.txt", "dest"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.transfer.calls, []string{"move a.txt"}; !slices.Equal(got, want) {
		t.Errorf("calls = %v, want %v", got, want)
	}

	if h.state.records[0].op != operation.Move || h.state.records[0].note != "note:move:file" {
		t.Errorf("recorded %+v, want the move note", h.state.records[0])
	}
}

// One argument names a source and no destination; it is refused instead of
// succeeding without doing anything.
func TestOneArgumentIsRefused(t *testing.T) {
	for _, op := range []operation.Operation{operation.Copy, operation.Move} {
		t.Run(string(op), func(t *testing.T) {
			h := newHarness("onlyone")

			err := h.run(transfercontent.Request{Op: op, Args: []string{"onlyone"}})
			if err == nil {
				t.Fatalf("%s with one argument returned no error", op)
			}

			if len(h.transfer.calls) != 0 || len(h.state.records) != 0 {
				t.Error("a refused transfer still did work")
			}
		})
	}
}

func TestEverySourceIsAttempted(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness("a.txt", "b.txt")
	h.transfer.err = boom
	h.transfer.errOn = "a.txt"

	err := h.run(transfercontent.Request{Op: operation.Copy, Args: []string{"a.txt", "b.txt", "dest"}})
	if !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}

	if got, want := h.transfer.calls, []string{"copy file a.txt", "copy file b.txt"}; !slices.Equal(got, want) {
		t.Errorf("calls = %v, want every source attempted", got)
	}

	if len(h.state.records) != 0 {
		t.Error("a partially failed transfer recorded an entry")
	}
}

func TestRecordErrorFailsTheRun(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness("a.txt")
	h.state.recordErr = boom

	if err := h.run(transfercontent.Request{Op: operation.Copy, Args: []string{"a.txt", "dest"}}); !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}
}
