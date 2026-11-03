package executedom

import (
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"

	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new execute note. params are the arguments the command ran with;
// without them the recorded line replays a different invocation.
func newExecuteCommandNote(command *command.Data, params []string) string {
	noteBuilder := note.New()

	noteBuilder.
		AddMainOperation(operations.Execute).
		AddSubOperation(operations.Command).
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

func NewExecuteLastNote(op operations.MainOperationEnum, command *command.Data, params []string) string {
	noteBuilder := note.New()

	noteBuilder.
		AddMainOperation(operations.Execute).
		AddSubOperation(operations.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec)

	if len(command.Steps) > 0 {
		noteBuilder.AddCommandSteps(command.Steps)
	}

	switch op {
	case operations.Execute:
		noteBuilder.AddCrowerExec("execute --last", params)
	case operations.Create:
		noteBuilder.AddCrowerExec("execute --create", params)
	case operations.Update:
		noteBuilder.AddCrowerExec("execute --update", params)
	}

	return noteBuilder.Build()
}
