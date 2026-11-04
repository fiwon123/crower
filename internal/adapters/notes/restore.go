package notes

import (
	"github.com/fiwon123/crower/internal/domain/operation"
)

// RestoreNote renders the history note for a restore.
//
// args are what the user passed (the command key when one was given); an
// interactive restore records the bare invocation.
func RestoreNote(args []string) string {
	return New().
		AddMainOperation(operation.Restore).
		AddCrowerExec("restore", args).
		Build()
}
