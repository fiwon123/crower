package createdom

import (
	"strings"
	"testing"

	command "github.com/fiwon123/crower/internal/data/command"
)

// The recorded crower_exec has to be the command that actually ran: the flag
// that selected the form, the name it created, the aliases it set, and each
// argument quoted so a value with a space stays one argument.
func TestCreateNotesRecordReplayableCommands(t *testing.T) {
	tests := []struct {
		name  string
		built string
		want  string
	}{
		{
			"a command create",
			newCreateCommandNote(command.New("greet", []string{"g"}, "echo hello world")),
			`crower_exec=crower create --alias 'g' 'greet' 'echo hello world'`,
		},
		{
			"a process create",
			newCreateProcessNote(command.New("editor", []string{"ed"}, "'/usr/bin/code'"), []string{"code runner"}),
			`crower_exec=crower create --process 'editor' --alias 'ed' 'code runner'`,
		},
		{
			"a script create",
			newCreateScriptCommandNote(command.New("deploy", nil, "deploy.sh")),
			`crower_exec=crower create --script 'deploy'`,
		},
		{
			"a reuse of the last execute",
			newCreateCommandLastExecuteNote(command.New("greet2", []string{"g2"}, "echo hi")),
			`crower_exec=crower create --execute --alias 'g2' 'greet2'`,
		},
		{
			"a file create keeps its flag",
			generateCreateFile([]string{"notes.md", "todo.txt"}),
			`crower_exec=crower create --file 'notes.md' 'todo.txt'`,
		},
		{
			"a folder create keeps its flag",
			newCreateFolder([]string{"tmp", "logs"}),
			`crower_exec=crower create --folder 'tmp' 'logs'`,
		},
		{
			"an alias with a comma survives the CSV flag",
			newCreateCommandNote(command.New("greet", []string{"a,b", "g2"}, "echo hi")),
			`crower_exec=crower create --alias '"a,b",g2' 'greet' 'echo hi'`,
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
