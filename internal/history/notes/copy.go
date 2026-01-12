package notes

import "github.com/fiwon123/crower/internal/data/state"

// Create a new copy note
func GenerateCopyNote(subOp state.SubOperationEnum, args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Copy).
		AddSubOperation(subOp).
		AddCrowerExec("crower copy", args).
		Build()
}
