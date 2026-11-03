package executedom

import (
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"

	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new execute note. params are the arguments the command ran with;
// without them the recorded line replays a different invocation.
func newExecuteCommandNote(command *command.Data, params []string) string {
	noteBuilder := note.New()

	noteBuilder.
		AddMainOperation(operation.Execute).
		AddSubOperation(operation.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec)

	// A command with steps never runs Exec; name the steps that did run.
	if len(command.Steps) > 0 {
		noteBuilder.AddCommandSteps(command.Steps)
	}

	return noteBuilder.
		AddCrowerExec("execute", append([]string{command.Name}, params...)).
		Build()
}

func NewExecuteLastNote(op operation.Operation, command *command.Data, params []string) string {
	noteBuilder := note.New()

	noteBuilder.
		AddMainOperation(operation.Execute).
		AddSubOperation(operation.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec)

	if len(command.Steps) > 0 {
		noteBuilder.AddCommandSteps(command.Steps)
	}

	switch op {
	case operation.Execute:
		noteBuilder.AddCrowerExec("execute --last", params)
	case operation.Create:
		noteBuilder.AddCrowerExec("execute --create", params)
	case operation.Update:
		noteBuilder.AddCrowerExec("execute --update", params)
	}

	return noteBuilder.Build()
}
