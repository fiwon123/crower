package upgradedom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new upgrade note
func newUpgradeNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Upgrade).
		AddCrowerExec("--upgrade", nil).
		Build()
}
