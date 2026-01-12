package notes

import "github.com/fiwon123/crower/internal/data/state"

// Create a new upgrade note
func GenerateUpgradeNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Upgrade).
		AddCrowerExec("--upgrade", nil).
		Build()
}
