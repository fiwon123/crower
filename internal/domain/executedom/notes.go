package executedom

import (
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"

	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new execute note
func newExecuteCommandNote(command *command.Data) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Execute).
		AddSubOperation(operations.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("execute", nil).
		Build()
}

func NewExecuteLastNote(op operations.MainOperationEnum, command *command.Data) string {
	noteBuilder := note.New()

	noteBuilder.
		AddMainOperation(operations.Execute).
		AddSubOperation(operations.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec)

	switch op {
	case operations.Execute:
		noteBuilder.AddCrowerExec("execute --last", nil)
	case operations.Create:
		noteBuilder.AddCrowerExec("execute --create", nil)
	case operations.Update:
		noteBuilder.AddCrowerExec("execute --update", nil)
	}

	return noteBuilder.Build()
}
