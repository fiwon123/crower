package copydom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new copy note
func newCopyNote(subOp operations.SubOperationEnum, args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Copy).
		AddSubOperation(subOp).
		AddCrowerExec("crower copy", args).
		Build()
}
