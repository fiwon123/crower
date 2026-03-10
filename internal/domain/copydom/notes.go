package copydom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new copy note
func newCopyNote(subOp operations.SubOperationEnum, args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Copy).
		AddSubOperation(subOp).
		AddCrowerExec("crower copy", args).
		Build()
}
