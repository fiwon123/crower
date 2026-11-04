package notes_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/shell"
)

// create --execute rebuilds a command from the note of the last execute. The
// reader looked for a bare "exec" key that this builder never wrote, so the two
// ends of the hand-off have to be tested against each other rather than each in
// isolation.
func TestExecuteCommandNoteCarriesTheExecCreateReuses(t *testing.T) {
	command := catalog.Command{Name: "greet", Aliases: []string{"g"}, Exec: "echo hello"}

	built := notes.ExecuteCommandNote(command, []string{"a b"})

	if got := notes.CommandExec(built); got != "echo hello" {
		t.Errorf("CommandExec(%q) = %q, want %q", built, got, "echo hello")
	}

	if got := notes.CommandName(built); got != "greet" {
		t.Errorf("CommandName(%q) = %q, want greet", built, got)
	}

	if want := "crower_exec=crower execute " + shell.Quote("greet") + " " + shell.Quote("a b"); !strings.Contains(built, want) {
		t.Errorf("note = %q, want it to contain %q", built, want)
	}

	dashed := notes.ExecuteCommandNote(command, []string{"-x"})
	if want := "crower_exec=crower execute -- " + shell.Quote("greet") + " " + shell.Quote("-x"); !strings.Contains(dashed, want) {
		t.Errorf("note = %q, want it to contain %q", dashed, want)
	}
}

// A command with steps never runs Exec; the note has to name the steps.
func TestExecuteNoteRecordsSteps(t *testing.T) {
	command := catalog.Command{Name: "multi", Exec: "unused", Steps: []string{"echo a", "echo b"}}

	built := notes.ExecuteCommandNote(command, nil)

	got, present := notes.CommandSteps(built)
	if !present || !slices.Equal(got, command.Steps) {
		t.Errorf("CommandSteps(%q) = %v, %v, want %v, true", built, got, present, command.Steps)
	}

	last := notes.ExecuteLastNote(operation.Execute, command, []string{"a"})
	if got, present := notes.CommandSteps(last); !present || !slices.Equal(got, command.Steps) {
		t.Errorf("CommandSteps(%q) = %v, %v, want the steps", last, got, present)
	}
}

// A declared shell is part of what ran; an empty one is the platform default,
// which every note written before the field existed already means.
func TestExecuteNoteRecordsTheShell(t *testing.T) {
	command := catalog.Command{Name: "deploy", Exec: "echo hi", Shell: "pwsh"}

	built := notes.ExecuteCommandNote(command, nil)
	if !strings.Contains(built, "command_shell=pwsh") {
		t.Errorf("note = %q, want it to record the shell", built)
	}

	last := notes.ExecuteLastNote(operation.Execute, command, nil)
	if !strings.Contains(last, "command_shell=pwsh") {
		t.Errorf("note = %q, want it to record the shell", last)
	}

	plain := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, nil)
	if strings.Contains(plain, "command_shell") {
		t.Errorf("note = %q, want no shell field for the platform default", plain)
	}
}

// Each --last form records the flag that reproduces it.
func TestExecuteLastNoteRecordsTheFlag(t *testing.T) {
	command := catalog.Command{Name: "greet", Exec: "echo hi"}

	tests := []struct {
		name string
		op   operation.Operation
		want string
	}{
		{"last", operation.Execute, "crower_exec=crower execute --last"},
		{"create", operation.Create, "crower_exec=crower execute --create"},
		{"update", operation.Update, "crower_exec=crower execute --update"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := notes.ExecuteLastNote(tt.op, command, []string{"p"})

			if want := tt.want + " " + shell.Quote("p"); !strings.Contains(got, want) {
				t.Errorf("note = %q, want it to contain %q", got, want)
			}
		})
	}
}
