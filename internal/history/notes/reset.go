package notes

import "github.com/fiwon123/crower/internal/data/state"

// Create a new reset note
func GenerateResetNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Reset).
		AddCrowerExec("reset", nil).
		Build()
}
