package executedom

import (
	"testing"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"
)

// create --execute rebuilds a command from the note of the last execute. The
// reader looked for a bare "exec" key that this builder never wrote, so the two
// ends of the hand-off have to be tested against each other rather than each in
// isolation.
func TestExecuteCommandNoteCarriesTheExecCreateReuses(t *testing.T) {
	cmd := command.New("greet", []string{"g"}, "echo hello")

	built := newExecuteCommandNote(cmd)

	if got := note.CommandExec(built); got != "echo hello" {
		t.Errorf("CommandExec(%q) = %q, want %q", built, got, "echo hello")
	}

	if got := note.CommandName(built); got != "greet" {
		t.Errorf("CommandName(%q) = %q, want greet", built, got)
	}
}
