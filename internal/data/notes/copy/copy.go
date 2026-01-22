package copynotesdata

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new copy note
func NewCopyNote(subOp operationsdata.SubOperationEnum, args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Copy).
		AddSubOperation(subOp).
		AddCrowerExec("crower copy", args).
		Build()
}
