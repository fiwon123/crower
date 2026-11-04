package notes

import (
	"github.com/fiwon123/crower/internal/domain/operation"
)

// ListCommandsNote renders the history note for `list`.
func ListCommandsNote() string {
	return New().
		AddMainOperation(operation.List).
		AddSubOperation(operation.Command).
		AddCrowerExec("list", nil).
		Build()
}

// ListProcessNote renders the history note for `list --process`.
func ListProcessNote(args []string) string {
	return New().
		AddMainOperation(operation.List).
		AddSubOperation(operation.Process).
		AddCrowerExec("list --process", args).
		Build()
}

// ListHistoryNote renders the history note for `list --history`.
func ListHistoryNote() string {
	return New().
		AddMainOperation(operation.List).
		AddSubOperation(operation.History).
		AddCrowerExec("list --history", nil).
		Build()
}

// ListFolderNote renders the history note for `list --folder`.
func ListFolderNote(args []string) string {
	return New().
		AddMainOperation(operation.List).
		AddSubOperation(operation.Folder).
		AddCrowerExec("list --folder", args).
		Build()
}

// ListSystemNote renders the history note for `list --system`.
func ListSystemNote() string {
	return New().
		AddMainOperation(operation.List).
		AddSubOperation(operation.System).
		AddCrowerExec("list --system", nil).
		Build()
}

// ListPathNote renders the history note for `list --syspath`.
func ListPathNote() string {
	return New().
		AddMainOperation(operation.List).
		AddSubOperation(operation.SysPath).
		AddCrowerExec("list --syspath", nil).
		Build()
}
