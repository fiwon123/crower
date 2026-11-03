package notes

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Reset renders the history note for a reset.
//
// A reset touches the whole command set rather than one command, so the note
// carries only the operation and the line that reproduces it.
func Reset() string {
	return note.New().
		AddMainOperation(operation.Reset).
		AddCrowerExec("reset", nil).
		Build()
}
