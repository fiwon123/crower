package checkdom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new check note
func newCheckNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Check).
		AddCrowerExec("--check", nil).
		Build()
}
