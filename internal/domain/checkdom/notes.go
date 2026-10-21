package checkdom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new check note
func newCheckNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Check).
		AddCrowerExec("crower --check", nil).
		Build()
}
