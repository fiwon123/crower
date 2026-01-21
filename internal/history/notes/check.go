package notes

import operationsdata "github.com/fiwon123/crower/internal/data/operations"

// Create a new check note
func GenerateCheckNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Check).
		AddCrowerExec("crower --check", nil).
		Build()
}
