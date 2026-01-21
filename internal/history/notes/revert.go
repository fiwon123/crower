package notes

import operationsdata "github.com/fiwon123/crower/internal/data/operations"

// Create a new revert note
func GenerateRevertNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Revert).
		AddCrowerExec("revert", args).
		Build()
}
