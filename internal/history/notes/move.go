package notes

import "github.com/fiwon123/crower/internal/data/state"

// Create a new move note
func GenerateMoveNote(subOp state.SubOperationEnum, args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Move).
		AddSubOperation(subOp).
		AddCrowerExec("crower move", args).
		Build()
}
