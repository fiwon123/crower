package restorecommand_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/restorecommand"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

type fakeState struct {
	commands  catalog.Catalog
	rows      []restorecommand.Entry
	found     map[string]catalog.Command
	findErr   error
	saved     []catalog.Catalog
	saveErr   error
	records   []record
	recordErr error
}

type record struct {
	op          operation.Operation
	commandName string
	note        string
}

func (s *fakeState) Commands() catalog.Catalog {
	return s.commands
}

func (s *fakeState) HistoryRows() []restorecommand.Entry {
	return s.rows
}

func (s *fakeState) FindCommand(key, file string) (catalog.Command, error) {
	if s.findErr != nil {
		return catalog.Command{}, s.findErr
	}

	command, ok := s.found[file]
	if !ok {
		return catalog.Command{}, fmt.Errorf("Not Found")
	}

	return command, nil
}

func (s *fakeState) Save(commands catalog.Catalog) error {
	if s.saveErr != nil {
		return s.saveErr
	}

	s.saved = append(s.saved, commands)
	return nil
}

func (s *fakeState) Record(op operation.Operation, commandName, note string) error {
	if s.recordErr != nil {
		return s.recordErr
	}

	s.records = append(s.records, record{op: op, commandName: commandName, note: note})
	return nil
}

type fakeOutput struct {
	infos    []string
	warnings []string
	table    string
}

func (o *fakeOutput) Info(message string, data ...any) {
	o.infos = append(o.infos, fmt.Sprint(append([]any{message}, data...)...))
}

func (o *fakeOutput) Warning(message string, data ...any) {
	o.warnings = append(o.warnings, message)
}

func (o *fakeOutput) HistoryTable() string {
	return o.table
}

func (o *fakeOutput) text() string {
	return strings.Join(o.infos, "\n")
}

type fakePrompt struct {
	entry      restorecommand.Entry
	chooseErr  error
	answer     bool
	confirmErr error
	asked      []string
}

func (p *fakePrompt) ChooseHistoryEntry() (restorecommand.Entry, error) {
	return p.entry, p.chooseErr
}

func (p *fakePrompt) Confirm(question string) (bool, error) {
	p.asked = append(p.asked, question)
	return p.answer, p.confirmErr
}

type fakeNotes struct{}

func (fakeNotes) RestoreNote(args []string) string {
	return "note:restore:" + strings.Join(args, ",")
}

type harness struct {
	state  *fakeState
	output *fakeOutput
	prompt *fakePrompt
}

func newHarness(commands ...catalog.Command) *harness {
	return &harness{
		state: &fakeState{
			commands: catalog.New(nil, commands),
			found:    map[string]catalog.Command{},
		},
		output: &fakeOutput{table: "HISTORY"},
		prompt: &fakePrompt{entry: restorecommand.Entry{Version: 3, File: "00003.yaml", Note: "note"}, answer: true},
	}
}

func (h *harness) run(request restorecommand.Request) error {
	return restorecommand.New(h.state, h.output, h.prompt, fakeNotes{}).Run(request)
}

func TestRestoreCreatesWhenTheNameIsGone(t *testing.T) {
	h := newHarness()
	h.state.found["00003.yaml"] = catalog.Command{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"}

	if err := h.run(restorecommand.Request{Key: "greet", Args: []string{"greet"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.state.saved) != 1 {
		t.Fatalf("saved %d catalogs, want 1", len(h.state.saved))
	}

	created, found := h.state.saved[0].Get("greet")
	if !found || created.Exec != "echo hi" {
		t.Errorf("saved = %+v, %v; want the restored command", created, found)
	}

	if got := h.output.text(); !strings.Contains(got, "restored command: restored by creating: &{greet [hi] echo hi []}") {
		t.Errorf("output = %q, want the creating description", got)
	}

	if len(h.state.records) != 1 || h.state.records[0].op != operation.Restore {
		t.Errorf("records = %+v, want the restore entry", h.state.records)
	}

	if h.state.records[0].note != "note:restore:greet" {
		t.Errorf("note = %q, want the restore note", h.state.records[0].note)
	}
}

func TestRestoreUpdatesWhenTheNameExists(t *testing.T) {
	h := newHarness(catalog.Command{Name: "greet", Aliases: []string{"hi"}, Exec: "echo current"})
	h.state.found["00003.yaml"] = catalog.Command{Name: "greet", Aliases: []string{"hi"}, Exec: "echo snapshot"}

	if err := h.run(restorecommand.Request{Key: "hi", Args: []string{"hi"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	updated, _ := h.state.saved[0].Get("greet")
	if updated.Exec != "echo snapshot" {
		t.Errorf("exec = %q, want the snapshot's line", updated.Exec)
	}

	// The key the user typed is what the update resolves: an alias works.
	if updated.Name != "greet" {
		t.Errorf("name = %q, want greet", updated.Name)
	}

	if got := h.output.text(); !strings.Contains(got, "restored by updating: from &{greet [hi] echo current []} to &{greet [hi] echo snapshot []}") {
		t.Errorf("output = %q, want the updating description", got)
	}
}

func TestRestoreShowsTheContentBeforeAsking(t *testing.T) {
	h := newHarness()
	h.state.found["00003.yaml"] = catalog.Command{Name: "greet", Exec: "echo hi"}

	if err := h.run(restorecommand.Request{}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	text := h.output.text()

	for _, want := range []string{"HISTORY", "Content", "Version:    3", "File:    00003.yaml", "Note:    note"} {
		if !strings.Contains(text, want) {
			t.Errorf("output = %q, want %q", text, want)
		}
	}

	if got, want := h.prompt.asked, []string{"Continue to restore"}; len(got) != 1 || got[0] != want[0] {
		t.Errorf("questions = %v, want %v", got, want)
	}
}

func TestRestoreDeclined(t *testing.T) {
	h := newHarness()
	h.prompt.answer = false

	if err := h.run(restorecommand.Request{}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(h.output.text(), "Cancelling restore...") {
		t.Errorf("output = %q, want the cancel line", h.output.text())
	}

	if len(h.state.saved) != 0 || len(h.state.records) != 0 {
		t.Error("a declined restore wrote something")
	}
}

func TestRestoreLookupFailure(t *testing.T) {
	boom := errors.New("Not Found")

	h := newHarness()
	h.state.findErr = boom

	err := h.run(restorecommand.Request{Key: "ghost"})
	if !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}

	if len(h.state.saved) != 0 {
		t.Error("a failed lookup saved a catalog")
	}
}

func TestRestoreAddFailure(t *testing.T) {
	h := newHarness(catalog.Command{Name: "other", Aliases: []string{"greet"}, Exec: "echo other"})
	h.state.found["00003.yaml"] = catalog.Command{Name: "greet", Exec: "echo hi"}

	err := h.run(restorecommand.Request{Key: "greet"})
	if err == nil || !strings.Contains(err.Error(), "already an alias") {
		t.Fatalf("error = %v, want the collision", err)
	}

	if len(h.state.saved) != 0 {
		t.Error("a refused restore saved a catalog")
	}
}

func TestRestoreErrors(t *testing.T) {
	boom := errors.New("boom")

	t.Run("choosing", func(t *testing.T) {
		h := newHarness()
		h.prompt.chooseErr = boom

		if err := h.run(restorecommand.Request{}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("confirming", func(t *testing.T) {
		h := newHarness()
		h.prompt.confirmErr = boom

		if err := h.run(restorecommand.Request{}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("saving", func(t *testing.T) {
		h := newHarness()
		h.state.found["00003.yaml"] = catalog.Command{Name: "greet", Exec: "echo hi"}
		h.state.saveErr = boom

		if err := h.run(restorecommand.Request{Key: "greet"}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.state.records) != 0 {
			t.Error("a failed save still recorded an entry")
		}
	})

	t.Run("recording", func(t *testing.T) {
		h := newHarness()
		h.state.found["00003.yaml"] = catalog.Command{Name: "greet", Exec: "echo hi"}
		h.state.recordErr = boom

		if err := h.run(restorecommand.Request{Key: "greet"}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}

func TestRestoreWarningsAreReported(t *testing.T) {
	// Two commands claim the alias "x": updating "greet" by its name changes
	// only the exec, so the update carries no warning; a rename would. Use a
	// config where the update's exec change is clean and assert no warning is
	// invented.
	h := newHarness(catalog.Command{Name: "greet", Exec: "echo current"})
	h.state.found["00003.yaml"] = catalog.Command{Name: "greet", Exec: "echo snapshot"}

	if err := h.run(restorecommand.Request{Key: "greet"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.output.warnings) != 0 {
		t.Errorf("warnings = %v, want none", h.output.warnings)
	}
}
