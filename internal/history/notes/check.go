package notes

import "github.com/fiwon123/crower/internal/data/state"

// Create a new check note
func GenerateCheckNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Check).
		AddCrowerExec("crower --check", nil).
		Build()
}
