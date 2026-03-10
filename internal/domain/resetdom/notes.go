package resetdom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new reset note
func newResetNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Reset).
		AddCrowerExec("reset", nil).
		Build()
}
