package listdom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new list note
func newListCommandsNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.List).
		AddSubOperation(operation.Command).
		AddCrowerExec("list", nil).
		Build()
}

func newListProcessNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.List).
		AddSubOperation(operation.Process).
		AddCrowerExec("list --process", args).
		Build()
}

func newListHistoriesNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.List).
		AddSubOperation(operation.History).
		AddCrowerExec("list --history", nil).
		Build()
}

func newListFolderNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.List).
		AddSubOperation(operation.Folder).
		AddCrowerExec("list --folder", args).
		Build()
}

func newListSystemNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.List).
		AddSubOperation(operation.System).
		AddCrowerExec("list --system", nil).
		Build()
}

func newListSystemPathNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.List).
		AddSubOperation(operation.SysPath).
		AddCrowerExec("list --syspath", nil).
		Build()
}
