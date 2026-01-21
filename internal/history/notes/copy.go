package notes

import operationsdata "github.com/fiwon123/crower/internal/data/operations"

// Create a new copy note
func GenerateCopyNote(subOp operationsdata.SubOperationEnum, args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Copy).
		AddSubOperation(subOp).
		AddCrowerExec("crower copy", args).
		Build()
}
