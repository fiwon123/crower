package notes

import "github.com/fiwon123/crower/internal/data/state"

// Create a new revert note
func GenerateRevertNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Revert).
		AddCrowerExec("revert", args).
		Build()
}
