package upgradedom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new upgrade note
func newUpgradeNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Upgrade).
		AddCrowerExec("--upgrade", nil).
		Build()
}
