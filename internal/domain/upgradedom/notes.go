package upgradedom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new upgrade note
func newUpgradeNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Upgrade).
		AddCrowerExec("--upgrade", nil).
		Build()
}
