package notes

import operationsdata "github.com/fiwon123/crower/internal/data/operations"

// Create a new upgrade note
func GenerateUpgradeNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Upgrade).
		AddCrowerExec("--upgrade", nil).
		Build()
}
