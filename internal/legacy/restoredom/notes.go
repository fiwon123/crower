package restoredom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new restore note. args are what the user passed (the command key
// when one was given); an interactive restore records the bare invocation.
func newRestoreNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Restore).
		AddCrowerExec("restore", args).
		Build()

}
