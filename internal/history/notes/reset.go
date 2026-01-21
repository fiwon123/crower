package notes

import operationsdata "github.com/fiwon123/crower/internal/data/operations"

// Create a new reset note
func GenerateResetNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Reset).
		AddCrowerExec("reset", nil).
		Build()
}
