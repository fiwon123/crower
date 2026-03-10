package revertdom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new revert note
func newRevertNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Revert).
		AddCrowerExec("revert", args).
		Build()
}
