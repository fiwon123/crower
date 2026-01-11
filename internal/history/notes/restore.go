package notes

import "github.com/fiwon123/crower/internal/data/state"

// Create a new restore note
func GenerateRestoreNote(msg string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Restore).
		AddCrowerExec("restore", nil).
		Build()

}
