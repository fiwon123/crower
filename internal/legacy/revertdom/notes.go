package revertdom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new revert note
func newRevertNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Revert).
		AddCrowerExec("revert", args).
		Build()
}
