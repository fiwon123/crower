package notes

import (
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/state"
)

// Create a new execute note
func GenerateExecuteCommandNote(command *command.Data) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Execute).
		AddSubOperation(state.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("execute", nil).
		Build()
}

func GenerateExecuteLastNote(op state.MainOperationEnum, command *command.Data) string {
	noteBuilder := New()

	noteBuilder.
		AddMainOperation(state.Execute).
		AddSubOperation(state.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec)

	switch op {
	case state.Execute:
		noteBuilder.AddCrowerExec("execute --last", nil)
	case state.Create:
		noteBuilder.AddCrowerExec("execute --create", nil)
	case state.Update:
		noteBuilder.AddCrowerExec("execute --update", nil)
	}

	return noteBuilder.Build()
}
