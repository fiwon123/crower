package notes

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// UpgradeNote renders the history note for a self-update.
func UpgradeNote() string {
	return note.New().
		AddMainOperation(operation.Upgrade).
		AddCrowerExec("--upgrade", nil).
		Build()
}
