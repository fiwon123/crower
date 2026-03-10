package resetdom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new reset note
func newResetNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Reset).
		AddCrowerExec("reset", nil).
		Build()
}
