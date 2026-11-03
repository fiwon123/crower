package resetcommands_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/resetcommands"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// fakeState records what the use case asked it to store.
type fakeState struct {
	commands  catalog.Catalog
	committed *commit
	commitErr error
}

type commit struct {
	commands catalog.Catalog
	op       operation.Operation
	replay   string
}

func (s *fakeState) Commands() catalog.Catalog {
	return s.commands
}

func (s *fakeState) Commit(commands catalog.Catalog, op operation.Operation, replay string) error {
	if s.commitErr != nil {
		return s.commitErr
	}

	s.committed = &commit{commands: commands, op: op, replay: replay}
	return nil
}

// fakeOutput keeps every diagnostic so a test can assert what the user saw.
type fakeOutput struct {
	lines []string
	table string
	rows  int
}

func (o *fakeOutput) Info(message string) {
	o.lines = append(o.lines, message)
}

func (o *fakeOutput) Table(commands catalog.Catalog) string {
	o.rows = commands.Len()
	return o.table
}

func (o *fakeOutput) text() string {
	return strings.Join(o.lines, "\n")
}

type fakePrompt struct {
	answer bool
	err    error
	asked  []string
}

func (p *fakePrompt) Confirm(question string) (bool, error) {
	p.asked = append(p.asked, question)
	return p.answer, p.err
}

func newUseCase(state *fakeState, output *fakeOutput, prompt *fakePrompt) *resetcommands.UseCase {
	return resetcommands.New(state, output, prompt)
}

func withCommands(names ...string) catalog.Catalog {
	commands := make([]catalog.Command, 0, len(names))

	for _, name := range names {
		commands = append(commands, catalog.Command{Name: name, Exec: "echo " + name})
	}

	return catalog.New(nil, commands)
}

// A yes empties the catalog and records the operation with the subcommand that
// reproduces it.
func TestRunConfirms(t *testing.T) {
	state := &fakeState{commands: withCommands("greet", "build")}
	output := &fakeOutput{table: "TABLE"}
	prompt := &fakePrompt{answer: true}

	if err := newUseCase(state, output, prompt).Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if state.committed == nil {
		t.Fatal("a confirmed reset committed nothing")
	}

	if got := state.committed.commands.Len(); got != 0 {
		t.Errorf("committed %d commands, want an empty catalog", got)
	}

	if state.committed.op != operation.Reset {
		t.Errorf("committed operation %q, want reset", state.committed.op)
	}

	if state.committed.replay != "reset" {
		t.Errorf("committed replay %q, want reset", state.committed.replay)
	}

	if got := prompt.asked; len(got) != 1 || got[0] != "Continue to reset" {
		t.Errorf("asked %v, want one Continue to reset question", got)
	}
}

// What will be erased is shown before the question: the table is a diagnostic
// and it describes the set the use case was given.
func TestRunShowsTheTableBeforeAsking(t *testing.T) {
	state := &fakeState{commands: withCommands("greet")}
	output := &fakeOutput{table: "TABLE"}
	prompt := &fakePrompt{answer: false}

	if err := newUseCase(state, output, prompt).Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(output.text(), "TABLE") {
		t.Errorf("the table was not printed: %q", output.text())
	}

	if output.rows != 1 {
		t.Errorf("the table was asked for %d rows, want 1", output.rows)
	}

	if !strings.Contains(output.text(), "All commands will be erased...") {
		t.Errorf("the warning is missing: %q", output.text())
	}
}

// Declining is not a failure: nothing happened, and that was the answer.
func TestRunDeclines(t *testing.T) {
	state := &fakeState{commands: withCommands("greet")}
	output := &fakeOutput{}
	prompt := &fakePrompt{answer: false}

	if err := newUseCase(state, output, prompt).Run(); err != nil {
		t.Fatalf("a declined reset returned an error: %v", err)
	}

	if state.committed != nil {
		t.Fatalf("a declined reset committed %+v", state.committed)
	}

	if !strings.Contains(output.text(), "Cancelling reset...") {
		t.Errorf("the decline was not reported: %q", output.text())
	}

	if strings.Contains(output.text(), "reset all commands") {
		t.Errorf("a declined reset reported that it ran: %q", output.text())
	}
}

// A prompt with no answer -- --no-input, or a closed stream -- is an error, and
// nothing is erased.
func TestRunPropagatesThePromptError(t *testing.T) {
	want := errors.New("input is required")
	state := &fakeState{commands: withCommands("greet")}
	output := &fakeOutput{}
	prompt := &fakePrompt{err: want}

	err := newUseCase(state, output, prompt).Run()

	if !errors.Is(err, want) {
		t.Fatalf("Run error = %v, want %v", err, want)
	}

	if state.committed != nil {
		t.Fatal("a failed prompt still committed")
	}
}

// A failed commit is reported, so the exit status says the work did not happen.
func TestRunPropagatesTheCommitError(t *testing.T) {
	want := errors.New("disk full")
	state := &fakeState{commands: withCommands("greet"), commitErr: want}
	output := &fakeOutput{}
	prompt := &fakePrompt{answer: true}

	if err := newUseCase(state, output, prompt).Run(); !errors.Is(err, want) {
		t.Fatalf("Run error = %v, want %v", err, want)
	}
}
