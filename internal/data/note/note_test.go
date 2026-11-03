package note

import "testing"

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
