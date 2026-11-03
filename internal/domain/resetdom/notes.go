package resetdom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new reset note
func newResetNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Reset).
		AddCrowerExec("reset", nil).
		Build()
}
