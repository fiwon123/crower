package movedom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new move note
func newMoveNote(subOp operations.SubOperationEnum, args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Move).
		AddSubOperation(subOp).
		AddCrowerExec("crower move", args).
		Build()
}
