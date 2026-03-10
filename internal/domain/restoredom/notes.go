package restoredom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new restore note
func newRestoreNote(msg string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Restore).
		AddCrowerExec("restore", nil).
		Build()

}
