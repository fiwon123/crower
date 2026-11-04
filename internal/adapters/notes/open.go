package notes

import (
	"github.com/fiwon123/crower/internal/domain/operation"
)

// OpenNote renders the history note for the positional `open`.
func OpenNote(args []string) string {
	return New().
		AddMainOperation(operation.Open).
		AddCrowerExec("open", args).
		Build()
}

// OpenFileNote renders the history note for `open --file`.
func OpenFileNote(args []string) string {
	return New().
		AddMainOperation(operation.Open).
		AddSubOperation(operation.File).
		AddCrowerExec("open --file", args).
		Build()
}

// OpenFolderNote renders the history note for `open --folder`.
func OpenFolderNote(args []string) string {
	return New().
		AddMainOperation(operation.Open).
		AddSubOperation(operation.Folder).
		AddCrowerExec("open --folder", args).
		Build()
}

// OpenSystemNote renders the history note for `open --system`.
func OpenSystemNote() string {
	return New().
		AddMainOperation(operation.Open).
		AddSubOperation(operation.System).
		AddCrowerExec("open --system", nil).
		Build()
}
