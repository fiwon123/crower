package notes

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/data/command"
)

// Create a new add note
func GenerateCreateNote(command *command.Data, args []string) string {
	var builder strings.Builder
	builder.WriteString("crower create ")
	for _, arg := range args {
		builder.WriteString(arg + " ")
	}

	return fmt.Sprintf("command_name=%s;command_alias=%v;command_exec=%s;crower_exec=%s", command.Name, command.AllAlias, command.Exec, builder.String())
}

// Create a new add process note
func GenerateCreateProcessNote(command *command.Data, args []string) string {
	var builder strings.Builder
	builder.WriteString("crower create ")
	for _, arg := range args {
		builder.WriteString(arg + " ")
	}

	return fmt.Sprintf("commandName=%s;commandAlias=%v;commandExec=%s;crowerExec=%s", command.Name, command.AllAlias, command.Exec, builder.String())
}

func GenerateCreateSystemVariableNote(args []string) string {
	var builder strings.Builder
	builder.WriteString("crower create --system ")
	for _, arg := range args {
		builder.WriteString(arg + " ")
	}

	return fmt.Sprintf("crowerExec=%s", builder.String())
}

func GenerateCreateSystemPathVariableNote(args []string) string {
	var builder strings.Builder
	builder.WriteString("crower create --syspath ")
	for _, arg := range args {
		builder.WriteString(arg + " ")
	}

	return fmt.Sprintf("crowerExec=%s", builder.String())
}
