package executedom

import (
	"slices"
	"strings"
	"testing"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/shell"
)

// create --execute rebuilds a command from the note of the last execute. The
// reader looked for a bare "exec" key that this builder never wrote, so the two
// ends of the hand-off have to be tested against each other rather than each in
// isolation.
func TestExecuteCommandNoteCarriesTheExecCreateReuses(t *testing.T) {
	cmd := command.New("greet", []string{"g"}, "echo hello")

	built := newExecuteCommandNote(cmd, []string{"a b"})

	if got := note.CommandExec(built); got != "echo hello" {
		t.Errorf("CommandExec(%q) = %q, want %q", built, got, "echo hello")
	}

	if got := note.CommandName(built); got != "greet" {
		t.Errorf("CommandName(%q) = %q, want greet", built, got)
	}

	if want := "crower_exec=crower execute " + shell.Quote("greet") + " " + shell.Quote("a b"); !strings.Contains(built, want) {
		t.Errorf("note = %q, want it to contain %q", built, want)
	}

	dashed := newExecuteCommandNote(cmd, []string{"-x"})
	if want := "crower_exec=crower execute -- " + shell.Quote("greet") + " " + shell.Quote("-x"); !strings.Contains(dashed, want) {
		t.Errorf("note = %q, want it to contain %q", dashed, want)
	}
}

// A command with steps never runs Exec; the note has to name the steps.
func TestExecuteNoteRecordsSteps(t *testing.T) {
	cmd := command.New("multi", nil, "unused")
	cmd.Steps = []string{"echo a", "echo b"}

	built := newExecuteCommandNote(cmd, nil)

	got, present := note.CommandSteps(built)
	if !present || !slices.Equal(got, cmd.Steps) {
		t.Errorf("CommandSteps(%q) = %v, %v, want %v, true", built, got, present, cmd.Steps)
	}

	last := NewExecuteLastNote(operations.Execute, cmd, []string{"a"})
	if got, present := note.CommandSteps(last); !present || !slices.Equal(got, cmd.Steps) {
		t.Errorf("CommandSteps(%q) = %v, %v, want the steps", last, got, present)
	}
}
