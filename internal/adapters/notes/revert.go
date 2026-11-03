package notes

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Revert renders the history note for a revert.
func Revert(args []string) string {
	return note.New().
		AddMainOperation(operation.Revert).
		AddCrowerExec("revert", args).
		Build()
}
