package notes

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// CopyNote renders the history note for a copy.
//
// subOp says what was copied -- file, folder or both -- so `list --history`
// reads as an audit line and not only as a replay recipe.
func CopyNote(subOp operation.SubOperation, args []string) string {
	return note.New().
		AddMainOperation(operation.Copy).
		AddSubOperation(subOp).
		AddCrowerExec("copy", args).
		Build()
}

// MoveNote renders the history note for a move.
func MoveNote(subOp operation.SubOperation, args []string) string {
	return note.New().
		AddMainOperation(operation.Move).
		AddSubOperation(subOp).
		AddCrowerExec("move", args).
		Build()
}
