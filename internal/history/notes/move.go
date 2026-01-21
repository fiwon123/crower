package notes

import operationsdata "github.com/fiwon123/crower/internal/data/operations"

// Create a new move note
func GenerateMoveNote(subOp operationsdata.SubOperationEnum, args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Move).
		AddSubOperation(subOp).
		AddCrowerExec("crower move", args).
		Build()
}
