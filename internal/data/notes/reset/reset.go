package resetnotesdata

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new reset note
func NewResetNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Reset).
		AddCrowerExec("reset", nil).
		Build()
}
