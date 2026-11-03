package upgradedom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new upgrade note
func newUpgradeNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Upgrade).
		AddCrowerExec("--upgrade", nil).
		Build()
}
