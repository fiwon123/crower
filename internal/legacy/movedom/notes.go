package movedom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new move note
func newMoveNote(subOp operation.SubOperation, args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Move).
		AddSubOperation(subOp).
		AddCrowerExec("move", args).
		Build()
}
