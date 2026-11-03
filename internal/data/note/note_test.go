package note

import (
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/shell"
)

// A value the format cannot carry raw has to come back exactly as it went in:
// the separator escape, the escape character itself, and both together.
func TestNoteValuesRoundTrip(t *testing.T) {
	exec := "make; echo done 100% %3B"

	built := New().
		AddMainOperation(operations.Execute).
		AddSubOperation(operations.Command).
		AddCommandName("build;all").
		AddCommandAlias([]string{"b;1", "b2"}).
		AddCommandExec(exec).
		AddCrowerExec("execute", []string{exec}).
		Build()

	if strings.Contains(built, "make; echo") {
		t.Errorf("the exec was written raw: %q", built)
	}

	if got := CommandName(built); got != "build;all" {
		t.Errorf("CommandName(%q) = %q, want build;all", built, got)
	}

	if got := CommandExec(built); got != exec {
		t.Errorf("CommandExec(%q) = %q, want %q", built, got, exec)
	}

	displayed := Display(built)
	if !strings.Contains(displayed, "command_name=build;all") || !strings.Contains(displayed, "command_exec="+exec) {
		t.Errorf("Display(%q) = %q, want the decoded values", built, displayed)
	}

	if !strings.Contains(displayed, "crower_exec=crower execute '"+exec+"'") {
		t.Errorf("Display(%q) = %q, want the quoted crower_exec decoded", built, displayed)
	}
}

func TestDisplay(t *testing.T) {
	tests := []struct {
		name  string
		built string
		want  string
	}{
		{"a plain note is unchanged", "main_operation=copy;crower_exec=crower copy a b", "main_operation=copy;crower_exec=crower copy a b"},
		{"separators are decoded", "command_exec=make%3B echo done", "command_exec=make; echo done"},
		{"percent is decoded", "command_exec=echo 100%25", "command_exec=echo 100%"},
		{"a field without a separator is kept", "bare", "bare"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Display(tt.built); got != tt.want {
				t.Errorf("Display(%q) = %q, want %q", tt.built, got, tt.want)
			}
		})
	}
}

func TestCommandExec(t *testing.T) {
	tests := []struct {
		name  string
		built string
		want  string
	}{
		{"an execute note", "main_operation=execute;sub_operation=command;command_name=greet;command_exec=echo hi", "echo hi"},
		{"no exec field", "main_operation=execute;sub_operation=command;command_name=greet", ""},
		{"a bare exec key is not the field", "main_operation=execute;sub_operation=command;exec=echo hi", ""},
		{"the old exec is not the current one", "main_operation=update;sub_operation=command;old_command_exec=old;command_exec=new", "new"},
		{"only an old exec", "main_operation=update;sub_operation=command;old_command_exec=old", ""},
		{"a file note carrying the text in a path", `main_operation=create;sub_operation=file;crower_exec=crower create --file "x;command_exec=evil"`, ""},
		{"an escaped separator", "main_operation=execute;sub_operation=command;command_exec=make%3B echo done", "make; echo done"},
		{"an old value with a percent", "main_operation=execute;sub_operation=command;command_exec=echo 100%", "echo 100%"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommandExec(tt.built); got != tt.want {
				t.Errorf("CommandExec(%q) = %q, want %q", tt.built, got, tt.want)
			}
		})
	}
}

func TestCommandName(t *testing.T) {
	tests := []struct {
		name  string
		built string
		want  string
	}{
		{"a create note", "main_operation=create;sub_operation=command;command_name=greet;command_exec=echo hi", "greet"},
		{"an execute note", "main_operation=execute;sub_operation=command;command_name=greet", "greet"},
		{"no command field", "main_operation=copy;crower_exec=crower copy a b", ""},
		{"old name is not the name", "main_operation=update;sub_operation=command;old_command_name=old;command_name=new", "new"},
		{"only an old name", "old_command_name=old", ""},
		{"a command note without a name", "main_operation=execute;sub_operation=command", ""},
		{"a file note carrying the text in a path", "main_operation=create;sub_operation=file;crower_exec=crower create --file \"x;command_name=greet\"", ""},
		{"an escaped separator", "main_operation=create;sub_operation=command;command_name=build%3Ball", "build;all"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommandName(tt.built); got != tt.want {
				t.Errorf("CommandName(%q) = %q, want %q", tt.built, got, tt.want)
			}
		})
	}
}

// A positional that looks like a flag needs a "--" before the positionals, or
// the recorded line fails to replay with "unknown flag".
func TestAddCrowerExecTerminator(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
		want    string
	}{
		{"no flag value", "copy", []string{"a", "b"}, "crower copy " + shell.Quote("a") + " " + shell.Quote("b")},
		{"a dash after another positional", "create", []string{"a", "-x"}, "crower create -- " + shell.Quote("a") + " " + shell.Quote("-x")},
		{"a literal terminator", "create", []string{"--", "a"}, "crower create -- " + shell.Quote("--") + " " + shell.Quote("a")},
		{"a lone dash", "create", []string{"-", "a"}, "crower create -- " + shell.Quote("-") + " " + shell.Quote("a")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			built := New().AddCrowerExec(tt.command, tt.args).Build()

			if want := "crower_exec=" + tt.want; !strings.Contains(built, want) {
				t.Errorf("note = %q, want it to contain %q", built, want)
			}
		})
	}
}

func TestCommandSteps(t *testing.T) {
	steps := []string{"echo a; b", "echo 100%"}

	built := New().
		AddMainOperation(operations.Execute).
		AddSubOperation(operations.Command).
		AddCommandSteps(steps).
		Build()

	got, present := CommandSteps(built)
	if !present || !slices.Equal(got, steps) {
		t.Errorf("CommandSteps(%q) = %v, %v, want %v, true", built, got, present, steps)
	}

	if _, present := CommandSteps("main_operation=execute;sub_operation=command"); present {
		t.Error("CommandSteps reported a field that is not there")
	}

	// A file note can carry the text inside a quoted path.
	fileNote := `main_operation=create;sub_operation=file;crower_exec=crower create --file "x;command_steps=[\"a\"]"`
	if _, present := CommandSteps(fileNote); present {
		t.Error("CommandSteps read a field out of a file note")
	}

	// A corrupt field is present, not absent: a caller must not fall back.
	if steps, present := CommandSteps("main_operation=execute;sub_operation=command;command_steps=not-json"); !present || steps != nil {
		t.Errorf("CommandSteps(corrupt) = %v, %v, want nil, true", steps, present)
	}
}

// A step with && must read as written, not HTML-escaped.
func TestDisplayKeepsStepsReadable(t *testing.T) {
	built := New().
		AddMainOperation(operations.Execute).
		AddSubOperation(operations.Command).
		AddCommandSteps([]string{"echo a && echo b"}).
		Build()

	if !strings.Contains(Display(built), `["echo a && echo b"]`) {
		t.Errorf("Display(%q) = %q, want the step as written", built, Display(built))
	}
}
