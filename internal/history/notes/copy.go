package notes

import "github.com/fiwon123/crower/internal/data/state"

// Create a new copy note
func GenerateCopyNote(subOp state.SubOperationEnum, args []string) string {
	noteBuilder := New()
	noteBuilder.AddMainOperation(state.Copy)
	noteBuilder.AddSubOperation(subOp)
	noteBuilder.AddCrowerExec("crower copy", args)

	return noteBuilder.Build()
}
