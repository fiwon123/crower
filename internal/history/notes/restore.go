package notes

import operationsdata "github.com/fiwon123/crower/internal/data/operations"

// Create a new restore note
func GenerateRestoreNote(msg string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Restore).
		AddCrowerExec("restore", nil).
		Build()

}
