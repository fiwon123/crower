package movedom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new move note
func newMoveNote(subOp operations.SubOperationEnum, args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Move).
		AddSubOperation(subOp).
		AddCrowerExec("crower move", args).
		Build()
}
