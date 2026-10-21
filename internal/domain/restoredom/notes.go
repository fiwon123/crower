package restoredom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new restore note
func newRestoreNote(msg string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Restore).
		AddCrowerExec("restore", nil).
		Build()

}
