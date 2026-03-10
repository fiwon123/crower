package movedom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new move note
func newMoveNote(subOp operationsdata.SubOperationEnum, args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Move).
		AddSubOperation(subOp).
		AddCrowerExec("crower move", args).
		Build()
}
