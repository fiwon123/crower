package updatecommand_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/updatecommand"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

type commit struct {
	commands    catalog.Catalog
	op          operation.Operation
	commandName string
	note        string
}

type fakeState struct {
	commands  catalog.Catalog
	last      map[operation.Operation]string
	committed *commit
	commitErr error
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

type fakeOutput struct {
	lines []string
	table string
	rows  int
}

func (o *fakeOutput) Info(message string, data ...any) {
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
	chosen  string
	name    string
	exec    string
	aliases []string
	answers map[string]bool
	err     error
	asked   []string
}

func (p *fakePrompt) record(question string) (bool, error) {
	p.asked = append(p.asked, question)

	if p.err != nil {
		return false, p.err
	}

	return p.answers[question], nil
}

func (p *fakePrompt) ChooseCommand() (string, error) {
	if p.err != nil {
		return "", p.err
	}

	p.asked = append(p.asked, "choose")
	return p.chosen, nil
}

func (p *fakePrompt) NewName() (string, error) {
	if p.err != nil {
		return "", p.err
	}

	p.asked = append(p.asked, "name")
	return p.name, nil
}

func (p *fakePrompt) NewAliases() ([]string, error) {
	if p.err != nil {
		return nil, p.err
	}

	p.asked = append(p.asked, "aliases")
	return p.aliases, nil
}

func (p *fakePrompt) NewExec() (string, error) {
	if p.err != nil {
		return "", p.err
	}

	p.asked = append(p.asked, "exec")
	return p.exec, nil
}

func (p *fakePrompt) Confirm(question string) (bool, error) {
	return p.record(question)
}

type fakeNotes struct {
	last operation.Operation
	got  bool
}

func (n *fakeNotes) Update(_, _ catalog.Command, last operation.Operation) string {
	n.last = last
	n.got = true
	return "update-note"
}

func catalogWith(commands ...catalog.Command) catalog.Catalog {
	return catalog.New(nil, commands)
}

func newUseCase(state *fakeState, output *fakeOutput, prompt *fakePrompt, notes *fakeNotes) *updatecommand.UseCase {
	return updatecommand.New(state, output, prompt, notes)
}

func greeter() *fakeState {
	return &fakeState{commands: catalogWith(catalog.Command{
		Name: "greet", Aliases: []string{"hi"}, Exec: "echo hello",
	})}
}

// A complete request still confirms once, then writes.
func TestRunAppliesASuppliedChange(t *testing.T) {
	state := greeter()
	output := &fakeOutput{}
	prompt := &fakePrompt{answers: map[string]bool{"Continue to update": true}}
	notes := &fakeNotes{}

	err := newUseCase(state, output, prompt, notes).Run(updatecommand.Request{
		Key:     "greet",
		Changes: catalog.Changes{Exec: "echo new"},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if state.committed == nil {
		t.Fatal("a confirmed update committed nothing")
	}

	if state.committed.op != operation.Update {
		t.Errorf("committed operation %q, want update", state.committed.op)
	}

	if state.committed.commandName != "greet" {
		t.Errorf("committed command %q, want greet", state.committed.commandName)
	}

	stored, _ := state.committed.commands.Get("greet")
	if stored.Exec != "echo new" {
		t.Errorf("committed exec = %q, want echo new", stored.Exec)
	}

	if !notes.got || notes.last != "" {
		t.Errorf("note rendered with last = %q, want the positional form", notes.last)
	}

	if !strings.Contains(output.text(), "updated command: ") {
		t.Errorf("the update was not reported: %q", output.text())
	}
}

// With no key, the table is shown and the row is asked for.
func TestRunAsksForTheCommand(t *testing.T) {
	state := greeter()
	output := &fakeOutput{table: "TABLE"}
	prompt := &fakePrompt{
		chosen:  "hi",
		answers: map[string]bool{"Continue to update": true},
	}
	notes := &fakeNotes{}

	err := newUseCase(state, output, prompt, notes).Run(updatecommand.Request{
		Changes: catalog.Changes{Exec: "echo new"},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(output.text(), "TABLE") {
		t.Errorf("the table was not shown: %q", output.text())
	}

	if output.rows != 1 {
		t.Errorf("the table was asked for %d rows, want 1", output.rows)
	}

	if state.committed == nil || state.committed.commandName != "greet" {
		t.Fatalf("committed %+v, want the command the alias resolves to", state.committed)
	}
}

// Every field the command line omitted is offered, and an accepted offer is
// asked for.
func TestRunAsksForMissingFields(t *testing.T) {
	state := greeter()
	output := &fakeOutput{}
	prompt := &fakePrompt{
		name:    "hello",
		aliases: []string{"hey"},
		exec:    "echo hi",
		answers: map[string]bool{
			"Do you want to update name":  true,
			"Do you want to update alias": true,
			"Do you want to update exec":  true,
			"Continue to update":          true,
		},
	}
	notes := &fakeNotes{}

	err := newUseCase(state, output, prompt, notes).Run(updatecommand.Request{Key: "greet"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	stored, found := state.committed.commands.Get("hello")
	if !found {
		t.Fatal("the command was not renamed")
	}

	if stored.Exec != "echo hi" {
		t.Errorf("exec = %q, want echo hi", stored.Exec)
	}

	if !strings.Contains(output.text(), "Old Command: ") || !strings.Contains(output.text(), "New Command: ") {
		t.Errorf("the before/after block is missing: %q", output.text())
	}

	switch strings.Join(prompt.asked, ",") {
	case "Do you want to update name,name,Do you want to update alias,aliases,Do you want to update exec,exec,Continue to update":
	default:
		t.Errorf("asked %v, want each offer then its value, then the confirmation", prompt.asked)
	}
}

// Declining is not a failure: nothing happened, and that was the answer.
func TestRunDeclines(t *testing.T) {
	state := greeter()
	output := &fakeOutput{}
	prompt := &fakePrompt{answers: map[string]bool{"Continue to update": false}}
	notes := &fakeNotes{}

	err := newUseCase(state, output, prompt, notes).Run(updatecommand.Request{
		Key:     "greet",
		Changes: catalog.Changes{Exec: "echo new"},
	})
	if err != nil {
		t.Fatalf("a declined update returned an error: %v", err)
	}

	if state.committed != nil {
		t.Fatalf("a declined update committed %+v", state.committed)
	}

	if !strings.Contains(output.text(), "Cancelling update...") {
		t.Errorf("the decline was not reported: %q", output.text())
	}
}

func TestRunReportsAnUnknownCommand(t *testing.T) {
	state := greeter()
	output := &fakeOutput{}
	prompt := &fakePrompt{answers: map[string]bool{"Continue to update": true}}
	notes := &fakeNotes{}

	if err := newUseCase(state, output, prompt, notes).Run(updatecommand.Request{Key: "nope"}); err == nil {
		t.Fatal("an unknown command returned no error")
	}

	if state.committed != nil {
		t.Fatal("an unknown command still committed")
	}
}

// A --last form takes no key and asks nothing: it would hang a scripted call.
func TestRunLastAsksNothing(t *testing.T) {
	state := greeter()
	state.last = map[operation.Operation]string{operation.Create: "greet"}
	output := &fakeOutput{}
	prompt := &fakePrompt{}
	notes := &fakeNotes{}

	err := newUseCase(state, output, prompt, notes).Run(updatecommand.Request{
		Last:    operation.Create,
		Changes: catalog.Changes{Exec: "echo new"},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(prompt.asked) != 0 {
		t.Errorf("a --last form asked %v, want no questions", prompt.asked)
	}

	if state.committed == nil || state.committed.op != operation.Update {
		t.Fatalf("committed %+v, want an update", state.committed)
	}

	if !notes.got || notes.last != operation.Create {
		t.Errorf("note rendered with last = %q, want create", notes.last)
	}
}

func TestRunLastWithoutAHistoryEntry(t *testing.T) {
	state := greeter()
	output := &fakeOutput{}
	prompt := &fakePrompt{}
	notes := &fakeNotes{}

	err := newUseCase(state, output, prompt, notes).Run(updatecommand.Request{
		Last:    operation.Execute,
		Changes: catalog.Changes{Exec: "echo new"},
	})

	if err == nil {
		t.Fatal("--last with no history returned no error")
	}

	if !strings.Contains(err.Error(), "execute") {
		t.Errorf("error = %v, want it to name the operation", err)
	}
}

func TestRunPropagatesTheCommitError(t *testing.T) {
	want := errors.New("disk full")
	state := greeter()
	state.commitErr = want
	output := &fakeOutput{}
	prompt := &fakePrompt{answers: map[string]bool{"Continue to update": true}}
	notes := &fakeNotes{}

	if err := newUseCase(state, output, prompt, notes).Run(updatecommand.Request{
		Key:     "greet",
		Changes: catalog.Changes{Exec: "echo new"},
	}); !errors.Is(err, want) {
		t.Fatalf("Run error = %v, want %v", err, want)
	}
}

func TestRunPropagatesAPromptError(t *testing.T) {
	want := errors.New("input is required")
	state := greeter()
	output := &fakeOutput{}
	prompt := &fakePrompt{err: want}
	notes := &fakeNotes{}

	if err := newUseCase(state, output, prompt, notes).Run(updatecommand.Request{
		Key:     "greet",
		Changes: catalog.Changes{Exec: "echo new"},
	}); !errors.Is(err, want) {
		t.Fatalf("Run error = %v, want %v", err, want)
	}

	if state.committed != nil {
		t.Fatal("a failed prompt still committed")
	}
}
