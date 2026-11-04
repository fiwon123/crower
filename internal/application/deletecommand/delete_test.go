package deletecommand_test

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/deletecommand"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// fakeState records what the use case asked it to store.
type fakeState struct {
	commands  catalog.Catalog
	last      map[operation.Operation]string
	committed *commit
	records   []record
	commitErr error
	recordErr error
}

type commit struct {
	commands    catalog.Catalog
	op          operation.Operation
	commandName string
	note        string
}

type record struct {
	op          operation.Operation
	commandName string
	note        string
}

func (s *fakeState) Commands() catalog.Catalog {
	return s.commands
}

func (s *fakeState) LastCommand(op operation.Operation) (string, bool) {
	name, ok := s.last[op]
	return name, ok
}

func (s *fakeState) Commit(commands catalog.Catalog, op operation.Operation, commandName, note string) error {
	if s.commitErr != nil {
		return s.commitErr
	}

	s.committed = &commit{commands: commands, op: op, commandName: commandName, note: note}
	return nil
}

func (s *fakeState) Record(op operation.Operation, commandName, note string) error {
	if s.recordErr != nil {
		return s.recordErr
	}

	s.records = append(s.records, record{op: op, commandName: commandName, note: note})
	return nil
}

// fakeOutput keeps every diagnostic so a test can assert what the user saw.
type fakeOutput struct {
	lines []string
	table string
	rows  int
}

func (o *fakeOutput) Info(message string, data ...any) {
	o.lines = append(o.lines, fmt.Sprint(append([]any{message}, data...)...))
}

func (o *fakeOutput) Table(commands catalog.Catalog) string {
	o.rows = commands.Len()
	return o.table
}

func (o *fakeOutput) text() string {
	return strings.Join(o.lines, "\n")
}

type fakePrompt struct {
	rows       []string
	chooseErr  error
	chooseCall int
	answer     bool
	confirmErr error
	asked      []string
}

func (p *fakePrompt) ChooseCommand() (string, error) {
	if p.chooseErr != nil {
		return "", p.chooseErr
	}

	p.chooseCall++

	if len(p.rows) == 0 {
		return "", nil
	}

	row := p.rows[0]
	p.rows = p.rows[1:]

	return row, nil
}

func (p *fakePrompt) Confirm(question string) (bool, error) {
	p.asked = append(p.asked, question)
	return p.answer, p.confirmErr
}

type fakeRemover struct {
	calls []string
	err   error
}

func (r *fakeRemover) File(path string) error {
	r.calls = append(r.calls, "file "+path)
	return r.err
}

func (r *fakeRemover) Folder(path string) error {
	r.calls = append(r.calls, "folder "+path)
	return r.err
}

type fakeEnvironment struct {
	calls []string
	out   string
	err   error
}

func (e *fakeEnvironment) Variable(name string) (string, error) {
	e.calls = append(e.calls, "variable "+name)
	return e.out, e.err
}

func (e *fakeEnvironment) VariablePath(path string) (string, error) {
	e.calls = append(e.calls, "path "+path)
	return e.out, e.err
}

type fakeHistory struct {
	table     string
	entry     deletecommand.Entry
	chooseErr error
	removeOut string
	removeErr error
	removed   []string
	chosen    bool
}

func (h *fakeHistory) HistoryTable() string {
	return h.table
}

func (h *fakeHistory) ChooseHistoryEntry() (deletecommand.Entry, error) {
	h.chosen = true
	return h.entry, h.chooseErr
}

func (h *fakeHistory) RemoveHistoryEntry(entry deletecommand.Entry) (string, error) {
	h.removed = append(h.removed, entry.File)
	return h.removeOut, h.removeErr
}

// fakeNotes renders notes that name the mode, so a test can tell which one the
// use case asked for without depending on the note format.
type fakeNotes struct {
	lastOp operation.Operation
}

func (n *fakeNotes) DeleteCommand(command catalog.Command, args []string) string {
	return "note:command:" + command.Name + ":" + fmt.Sprint(args)
}

func (n *fakeNotes) DeleteLast(op operation.Operation, command catalog.Command) string {
	n.lastOp = op
	return "note:last:" + string(op) + ":" + command.Name
}

func (n *fakeNotes) DeleteFile(args []string) string {
	return "note:file:" + fmt.Sprint(args)
}

func (n *fakeNotes) DeleteFolder(args []string) string {
	return "note:folder:" + fmt.Sprint(args)
}

func (n *fakeNotes) DeleteSystemVariable(args []string) string {
	return "note:system:" + fmt.Sprint(args)
}

func (n *fakeNotes) DeleteSystemPath(args []string) string {
	return "note:syspath:" + fmt.Sprint(args)
}

type harness struct {
	state   *fakeState
	output  *fakeOutput
	prompt  *fakePrompt
	remover *fakeRemover
	environ *fakeEnvironment
	history *fakeHistory
	notes   *fakeNotes
}

func newHarness(commands ...catalog.Command) *harness {
	return &harness{
		state: &fakeState{
			commands: catalog.New(nil, commands),
			last:     map[operation.Operation]string{},
		},
		output:  &fakeOutput{table: "TABLE"},
		prompt:  &fakePrompt{answer: true},
		remover: &fakeRemover{},
		environ: &fakeEnvironment{},
		history: &fakeHistory{},
		notes:   &fakeNotes{},
	}
}

func (h *harness) run(request deletecommand.Request) error {
	return deletecommand.New(
		h.state, h.output, h.prompt, h.remover, h.environ, h.history, h.notes,
	).Run(request)
}

func twoCommands() []catalog.Command {
	return []catalog.Command{
		{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"},
		{Name: "build", Exec: "go build ./..."},
	}
}

func TestDeleteCommandByName(t *testing.T) {
	h := newHarness(twoCommands()...)

	if err := h.run(deletecommand.Request{Mode: deletecommand.ModeCommand, Key: "greet"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.state.committed == nil {
		t.Fatal("a confirmed delete committed nothing")
	}

	if got, want := h.state.committed.commands.Names(), []string{"build"}; !slices.Equal(got, want) {
		t.Errorf("committed catalog = %v, want %v", got, want)
	}

	if h.state.committed.op != operation.Delete {
		t.Errorf("committed operation %q, want delete", h.state.committed.op)
	}

	if h.state.committed.commandName != "greet" {
		t.Errorf("committed command name %q, want greet", h.state.committed.commandName)
	}

	if want := "note:command:greet:[]"; h.state.committed.note != want {
		t.Errorf("committed note = %q, want %q", h.state.committed.note, want)
	}

	// A key was given, so nothing had to be selected.
	if h.output.rows != 0 {
		t.Errorf("the table was rendered for a keyed delete")
	}

	if h.prompt.chooseCall != 0 {
		t.Errorf("the use case asked for a row %d time(s), want 0", h.prompt.chooseCall)
	}

	if got, want := h.prompt.asked, []string{"Continue to delete"}; !slices.Equal(got, want) {
		t.Errorf("questions = %v, want %v", got, want)
	}

	if !strings.Contains(h.output.text(), "deleted command: greet") {
		t.Errorf("output = %q, want the success line", h.output.text())
	}

	if !strings.Contains(h.output.text(), "greet") || h.output.rows != 0 {
		t.Errorf("output = %q, want the command shown before the question", h.output.text())
	}
}

// The confirmation block shows the command the key resolved to, aliases
// included, whichever form the key took.
func TestDeleteCommandShowsTheResolvedCommand(t *testing.T) {
	h := newHarness(twoCommands()...)

	if err := h.run(deletecommand.Request{Mode: deletecommand.ModeCommand, Key: "hi"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	text := h.output.text()

	for _, want := range []string{"Name:    greet", "Aliases: [hi]", "Exec:    echo hi"} {
		if !strings.Contains(text, want) {
			t.Errorf("output = %q, want it to contain %q", text, want)
		}
	}

	if h.state.committed.commandName != "greet" {
		t.Errorf("committed command name %q, want greet", h.state.committed.commandName)
	}
}

func TestDeleteCommandAsksWhenNoKeyIsGiven(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.prompt.rows = []string{"build"}

	if err := h.run(deletecommand.Request{Mode: deletecommand.ModeCommand}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.output.rows != 2 {
		t.Errorf("table shows %d rows, want 2", h.output.rows)
	}

	if h.prompt.chooseCall != 1 {
		t.Errorf("ChooseCommand called %d time(s), want 1", h.prompt.chooseCall)
	}

	if h.state.committed == nil || h.state.committed.commandName != "build" {
		t.Errorf("committed = %+v, want build deleted", h.state.committed)
	}
}

func TestDeleteCommandUnknownKeyAsksAgain(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.prompt.rows = []string{"build"}

	if err := h.run(deletecommand.Request{Mode: deletecommand.ModeCommand, Key: "nope"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(h.output.text(), "Command not found, try to select one.") {
		t.Errorf("output = %q, want the retry notice", h.output.text())
	}

	if h.prompt.chooseCall != 1 {
		t.Errorf("ChooseCommand called %d time(s), want 1", h.prompt.chooseCall)
	}

	if h.state.committed == nil || h.state.committed.commandName != "build" {
		t.Errorf("committed = %+v, want build deleted", h.state.committed)
	}
}

func TestDeleteCommandUnknownKeyTwiceFails(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.prompt.rows = []string{"also-nope"}

	err := h.run(deletecommand.Request{Mode: deletecommand.ModeCommand, Key: "nope"})

	if !strings.Contains(err.Error(), "command not found") {
		t.Fatalf("error = %v, want CommandNotFound", err)
	}

	if !strings.Contains(err.Error(), "also-nope") {
		t.Errorf("error = %v, want it to name the second key", err)
	}

	if h.state.committed != nil {
		t.Error("a failed delete committed a catalog")
	}
}

func TestDeleteCommandDeclined(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.prompt.answer = false

	if err := h.run(deletecommand.Request{Mode: deletecommand.ModeCommand, Key: "greet"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.state.committed != nil {
		t.Error("a declined delete committed a catalog")
	}

	if !strings.Contains(h.output.text(), "Cancelling delete...") {
		t.Errorf("output = %q, want the cancel line", h.output.text())
	}
}

func TestDeleteCommandPromptErrors(t *testing.T) {
	boom := errors.New("boom")

	t.Run("selecting", func(t *testing.T) {
		h := newHarness(twoCommands()...)
		h.prompt.chooseErr = boom

		if err := h.run(deletecommand.Request{Mode: deletecommand.ModeCommand}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if h.state.committed != nil {
			t.Error("a prompt error committed a catalog")
		}
	})

	t.Run("confirming", func(t *testing.T) {
		h := newHarness(twoCommands()...)
		h.prompt.confirmErr = boom

		if err := h.run(deletecommand.Request{Mode: deletecommand.ModeCommand, Key: "greet"}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if h.state.committed != nil {
			t.Error("a prompt error committed a catalog")
		}
	})
}

func TestDeleteCommandCommitError(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness(twoCommands()...)
	h.state.commitErr = boom

	if err := h.run(deletecommand.Request{Mode: deletecommand.ModeCommand, Key: "greet"}); !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}

	if strings.Contains(h.output.text(), "deleted command:") {
		t.Errorf("output = %q, want no success line after a failed commit", h.output.text())
	}
}

// The --last forms take no key and ask nothing: the flags describe the whole
// request, and a prompt in the middle of a scripted call would hang it.
func TestDeleteLast(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.state.last[operation.Create] = "greet"

	if err := h.run(deletecommand.Request{Mode: deletecommand.ModeLast, Last: operation.Create}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.prompt.chooseCall != 0 || len(h.prompt.asked) != 0 {
		t.Errorf("the --last form asked %d row(s) and %d question(s), want none", h.prompt.chooseCall, len(h.prompt.asked))
	}

	if h.output.rows != 0 {
		t.Error("the --last form rendered the table")
	}

	if h.state.committed == nil || h.state.committed.commandName != "greet" {
		t.Fatalf("committed = %+v, want greet deleted", h.state.committed)
	}

	if want := "note:last:create:greet"; h.state.committed.note != want {
		t.Errorf("committed note = %q, want %q", h.state.committed.note, want)
	}

	if h.notes.lastOp != operation.Create {
		t.Errorf("Notes saw operation %q, want create", h.notes.lastOp)
	}
}

func TestDeleteLastWithoutAHistoryEntryFails(t *testing.T) {
	h := newHarness(twoCommands()...)

	err := h.run(deletecommand.Request{Mode: deletecommand.ModeLast, Last: operation.Execute})

	if !strings.Contains(err.Error(), "command not found") {
		t.Fatalf("error = %v, want CommandNotFound", err)
	}

	if !strings.Contains(err.Error(), string(operation.Execute)) {
		t.Errorf("error = %v, want it to name execute", err)
	}

	if h.state.committed != nil {
		t.Error("a --last delete with no history entry committed a catalog")
	}
}

func TestDeleteFileAndFolder(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		h := newHarness()

		if err := h.run(deletecommand.Request{Mode: deletecommand.ModeFile, Args: []string{"notes.txt"}}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if got, want := h.remover.calls, []string{"file notes.txt"}; !slices.Equal(got, want) {
			t.Errorf("remover calls = %v, want %v", got, want)
		}

		if want := "note:file:[notes.txt]"; len(h.state.records) != 1 || h.state.records[0].note != want {
			t.Errorf("records = %+v, want one with %q", h.state.records, want)
		}

		if h.state.committed != nil {
			t.Error("a file delete committed a command catalog")
		}

		assertRecordedWithoutCommand(t, h.state.records[0])
	})

	t.Run("folder", func(t *testing.T) {
		h := newHarness()

		if err := h.run(deletecommand.Request{Mode: deletecommand.ModeFolder, Args: []string{"reports"}}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if got, want := h.remover.calls, []string{"folder reports"}; !slices.Equal(got, want) {
			t.Errorf("remover calls = %v, want %v", got, want)
		}

		if want := "note:folder:[reports]"; len(h.state.records) != 1 || h.state.records[0].note != want {
			t.Errorf("records = %+v, want one with %q", h.state.records, want)
		}

		assertRecordedWithoutCommand(t, h.state.records[0])
	})
}

func TestDeletePathWithoutAnArgumentFails(t *testing.T) {
	tests := []struct {
		name string
		mode deletecommand.Mode
		want string
	}{
		{"file", deletecommand.ModeFile, "needs a file or folder path"},
		{"folder", deletecommand.ModeFolder, "needs a file or folder path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHarness()

			if err := h.run(deletecommand.Request{Mode: tt.mode}); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}

			if len(h.remover.calls) != 0 || len(h.state.records) != 0 {
				t.Errorf("remover calls = %v, records = %v; want neither", h.remover.calls, h.state.records)
			}
		})
	}
}

func TestDeleteFileRemoverError(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness()
	h.remover.err = boom

	if err := h.run(deletecommand.Request{Mode: deletecommand.ModeFile, Args: []string{"notes.txt"}}); !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}

	if len(h.state.records) != 0 {
		t.Errorf("records = %v, want none after a failed delete", h.state.records)
	}
}

func TestDeleteSystemVariableAndPath(t *testing.T) {
	t.Run("variable", func(t *testing.T) {
		h := newHarness()
		h.environ.out = "var name deleted"

		if err := h.run(deletecommand.Request{Mode: deletecommand.ModeSystemVariable, Args: []string{"FOO"}}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if got, want := h.environ.calls, []string{"variable FOO"}; !slices.Equal(got, want) {
			t.Errorf("environment calls = %v, want %v", got, want)
		}

		if !strings.Contains(h.output.text(), "var name deleted") {
			t.Errorf("output = %q, want the adapter's message", h.output.text())
		}

		if want := "note:system:[FOO]"; len(h.state.records) != 1 || h.state.records[0].note != want {
			t.Errorf("records = %+v, want one with %q", h.state.records, want)
		}

		assertRecordedWithoutCommand(t, h.state.records[0])
	})

	t.Run("path", func(t *testing.T) {
		h := newHarness()
		h.environ.out = "Deleted path from PATH"

		if err := h.run(deletecommand.Request{Mode: deletecommand.ModeSystemPath, Args: []string{"/opt/bin"}}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if got, want := h.environ.calls, []string{"path /opt/bin"}; !slices.Equal(got, want) {
			t.Errorf("environment calls = %v, want %v", got, want)
		}

		if want := "note:syspath:[/opt/bin]"; len(h.state.records) != 1 || h.state.records[0].note != want {
			t.Errorf("records = %+v, want one with %q", h.state.records, want)
		}

		assertRecordedWithoutCommand(t, h.state.records[0])
	})
}

func TestDeleteSystemWithoutAnArgumentFails(t *testing.T) {
	h := newHarness()

	err := h.run(deletecommand.Request{Mode: deletecommand.ModeSystemVariable})
	if err == nil || !strings.Contains(err.Error(), "missing arguments: var name") {
		t.Fatalf("error = %v, want NotEnoughArgs", err)
	}

	if len(h.environ.calls) != 0 || len(h.state.records) != 0 {
		t.Errorf("environment calls = %v, records = %v; want neither", h.environ.calls, h.state.records)
	}
}

func TestDeleteHistory(t *testing.T) {
	h := newHarness()
	h.history.table = "HISTORY"
	h.history.entry = deletecommand.Entry{Version: 7, File: "00007.yaml", Note: "greet"}
	h.history.removeOut = "history entry deleted"

	if err := h.run(deletecommand.Request{Mode: deletecommand.ModeHistory}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	text := h.output.text()

	for _, want := range []string{"HISTORY", "Content", "Version:    7", "File:    00007.yaml", "Note:    greet", "history entry deleted"} {
		if !strings.Contains(text, want) {
			t.Errorf("output = %q, want it to contain %q", text, want)
		}
	}

	if got, want := h.history.removed, []string{"00007.yaml"}; !slices.Equal(got, want) {
		t.Errorf("removed = %v, want %v", got, want)
	}

	if h.prompt.asked == nil || h.prompt.asked[0] != "Continue to delete history content" {
		t.Errorf("questions = %v, want the history confirmation", h.prompt.asked)
	}
}

func TestDeleteHistoryDeclined(t *testing.T) {
	h := newHarness()
	h.prompt.answer = false

	if err := h.run(deletecommand.Request{Mode: deletecommand.ModeHistory}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.history.removed) != 0 {
		t.Errorf("removed = %v, want nothing", h.history.removed)
	}

	if !strings.Contains(h.output.text(), "Cancelling Delete History Content...") {
		t.Errorf("output = %q, want the cancel line", h.output.text())
	}
}

// A history entry that cannot be chosen or removed fails the run and records
// nothing -- in particular it does not fall back to a command delete.
func TestDeleteHistoryErrors(t *testing.T) {
	boom := errors.New("boom")

	t.Run("choosing", func(t *testing.T) {
		h := newHarness()
		h.history.chooseErr = boom

		if err := h.run(deletecommand.Request{Mode: deletecommand.ModeHistory}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.history.removed) != 0 {
			t.Error("a failed choose still removed an entry")
		}
	})

	t.Run("removing", func(t *testing.T) {
		h := newHarness()
		h.history.removeErr = boom

		if err := h.run(deletecommand.Request{Mode: deletecommand.ModeHistory}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}

// A zero Request is the positional form, so an empty Key asks for a row rather
// than failing.
func TestZeroRequestAsksForACommand(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.prompt.rows = []string{"greet"}

	if err := h.run(deletecommand.Request{}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.prompt.chooseCall != 1 {
		t.Errorf("ChooseCommand called %d time(s), want 1", h.prompt.chooseCall)
	}
}

func assertRecordedWithoutCommand(t *testing.T, got record) {
	t.Helper()

	if got.op != operation.Delete {
		t.Errorf("recorded operation %q, want delete", got.op)
	}

	if got.commandName != "" {
		t.Errorf("recorded command name %q, want empty", got.commandName)
	}
}
