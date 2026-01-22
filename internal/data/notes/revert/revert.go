package revertnotesdata

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new revert note
func NewRevertNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Revert).
		AddCrowerExec("revert", args).
		Build()
}
