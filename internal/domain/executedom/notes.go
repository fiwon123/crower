package executedom

import (
	commanddata "github.com/fiwon123/crower/internal/data/command"
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new execute note
func newExecuteCommandNote(command *commanddata.Data) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Execute).
		AddSubOperation(operationsdata.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("execute", nil).
		Build()
}

func NewExecuteLastNote(op operationsdata.MainOperationEnum, command *commanddata.Data) string {
	noteBuilder := notesdata.New()

	noteBuilder.
		AddMainOperation(operationsdata.Execute).
		AddSubOperation(operationsdata.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec)

	switch op {
	case operationsdata.Execute:
		noteBuilder.AddCrowerExec("execute --last", nil)
	case operationsdata.Create:
		noteBuilder.AddCrowerExec("execute --create", nil)
	case operationsdata.Update:
		noteBuilder.AddCrowerExec("execute --update", nil)
	}

	return noteBuilder.Build()
}
