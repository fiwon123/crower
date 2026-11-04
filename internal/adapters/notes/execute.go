package notes

import (
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// ExecuteCommandNote renders the history note for `execute <command>`.
//
// params are the arguments the command ran with; without them the recorded line
// replays a different invocation.
func ExecuteCommandNote(command catalog.Command, params []string) string {
	builder := New()

	builder.
		AddMainOperation(operation.Execute).
		AddSubOperation(operation.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.Aliases).
		AddCommandExec(command.Exec)

	// A command with steps never runs Exec; name the steps that did run.
	if len(command.Steps) > 0 {
		builder.AddCommandSteps(command.Steps)
	}

	// An empty shell is the platform default, which is what every note meant
	// before the field existed, so only a declared shell is recorded.
	if command.Shell != "" {
		builder.AddCommandShell(command.Shell)
	}

	return builder.
		AddCrowerExec("execute", append([]string{command.Name}, params...)).
		Build()
}

// ExecuteLastNote renders the history note for the --last/--create/--update
// forms.
func ExecuteLastNote(op operation.Operation, command catalog.Command, params []string) string {
	builder := New()

	builder.
		AddMainOperation(operation.Execute).
		AddSubOperation(operation.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.Aliases).
		AddCommandExec(command.Exec)

	if len(command.Steps) > 0 {
		builder.AddCommandSteps(command.Steps)
	}

	if command.Shell != "" {
		builder.AddCommandShell(command.Shell)
	}

	switch op {
	case operation.Execute:
		builder.AddCrowerExec("execute --last", params)
	case operation.Create:
		builder.AddCrowerExec("execute --create", params)
	case operation.Update:
		builder.AddCrowerExec("execute --update", params)
	}

	return builder.Build()
}
