package executedom

import (
	commanddata "github.com/fiwon123/crower/internal/data/command"
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new execute note
func newExecuteCommandNote(command *commanddata.Data) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Execute).
		AddSubOperation(operations.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("execute", nil).
		Build()
}

func NewExecuteLastNote(op operations.MainOperationEnum, command *commanddata.Data) string {
	noteBuilder := notesdata.New()

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
