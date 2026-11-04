package notes

import (
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Check renders the history note for `--check`.
func Check() string {
	return New().
		AddMainOperation(operation.Check).
		AddCrowerExec("--check", nil).
		Build()
}
