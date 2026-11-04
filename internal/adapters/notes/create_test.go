package notes_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/shell"
)

// The recorded crower_exec has to be the command that actually ran: the flag
// that selected the form, the name it created, the aliases it set, and each
// argument quoted so a value with a space stays one argument.
func TestCreateNotesRecordReplayableCommands(t *testing.T) {
	// Expected values follow the recording OS, so the test runs on every
	// platform the release ships.
	quote := shell.Quote

	tests := []struct {
		name  string
		built string
		want  string
	}{
		{
			"a command create",
			notes.CreateCommandNote(catalog.Command{Name: "greet", Aliases: []string{"g"}, Exec: "echo hello world"}),
			"crower_exec=crower create --alias " + quote("g") + " " + quote("greet") + " " + quote("echo hello world"),
		},
		{
			"a process create",
			notes.CreateProcessNote(catalog.Command{Name: "editor", Aliases: []string{"ed"}, Exec: "'/usr/bin/code'"}, []string{"code runner"}),
			"crower_exec=crower create --process " + quote("editor") + " --alias " + quote("ed") + " " + quote("code runner"),
		},
		{
			"a script create",
			notes.CreateScriptNote(catalog.Command{Name: "deploy", Exec: "deploy.sh"}),
			"crower_exec=crower create --script " + quote("deploy"),
		},
		{
			"a reuse of the last execute",
			notes.CreateLastNote(catalog.Command{Name: "greet2", Aliases: []string{"g2"}, Exec: "echo hi"}),
			"crower_exec=crower create --execute --alias " + quote("g2") + " " + quote("greet2"),
		},
		{
			"a file create keeps its flag",
			notes.CreateFileNote([]string{"notes.md", "todo.txt"}),
			"crower_exec=crower create --file " + quote("notes.md") + " " + quote("todo.txt"),
		},
		{
			"a folder create keeps its flag",
			notes.CreateFolderNote([]string{"tmp", "logs"}),
			"crower_exec=crower create --folder " + quote("tmp") + " " + quote("logs"),
		},
		{
			"a system variable create keeps its flag",
			notes.CreateSystemVariableNote([]string{"FOO", "bar"}),
			"crower_exec=crower create --system " + quote("FOO") + " " + quote("bar"),
		},
		{
			"a path create keeps its flag",
			notes.CreateSystemPathNote([]string{"/opt/tools"}),
			"crower_exec=crower create --syspath " + quote("/opt/tools"),
		},
		{
			"an alias with a comma survives the CSV flag",
			notes.CreateCommandNote(catalog.Command{Name: "greet", Aliases: []string{"a,b", "g2"}, Exec: "echo hi"}),
			"crower_exec=crower create --alias " + quote(`"a,b",g2`) + " " + quote("greet") + " " + quote("echo hi"),
		},
		{
			"a name that looks like a flag needs the terminator",
			notes.CreateCommandNote(catalog.Command{Name: "-v", Exec: "echo x"}),
			"crower_exec=crower create -- " + quote("-v") + " " + quote("echo x"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.built, tt.want) {
				t.Errorf("note = %q, want it to contain %q", tt.built, tt.want)
			}
		})
	}
}

// The note also records what the command is, so `list --history` reads as an
// audit line and not only as a replay recipe.
func TestCreateCommandNoteRecordsTheCommand(t *testing.T) {
	got := notes.CreateCommandNote(catalog.Command{Name: "greet", Aliases: []string{"g"}, Exec: "echo hi"})

	for _, want := range []string{
		"main_operation=create",
		"sub_operation=command",
		"command_name=greet",
		"command_alias=[g]",
		"command_exec=echo hi",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("note = %q, want it to contain %q", got, want)
		}
	}
}
