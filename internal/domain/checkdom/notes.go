package checkdom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new check note
func newCheckNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Check).
		AddCrowerExec("crower --check", nil).
		Build()
}
