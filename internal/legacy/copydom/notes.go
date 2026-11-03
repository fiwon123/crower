package copydom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new copy note
func newCopyNote(subOp operation.SubOperation, args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Copy).
		AddSubOperation(subOp).
		AddCrowerExec("copy", args).
		Build()
}
