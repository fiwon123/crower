package notes

import (
	"github.com/fiwon123/crower/internal/data/state"
)

// Create a new list note
func GenerateListCommandsNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.List).
		AddSubOperation(state.Command).
		AddCrowerExec("list", nil).
		Build()
}

func GenerateListProcessNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.List).
		AddSubOperation(state.Process).
		AddCrowerExec("list --process", args).
		Build()
}

func GenerateListHistoriesNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.List).
		AddSubOperation(state.History).
		AddCrowerExec("list --history", nil).
		Build()
}

func GenerateListFolderNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.List).
		AddSubOperation(state.Folder).
		AddCrowerExec("list --folder", args).
		Build()
}

func GenerateListSystemNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.List).
		AddSubOperation(state.System).
		AddCrowerExec("list --system", nil).
		Build()
}

func GenerateListSystemPathNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.List).
		AddSubOperation(state.SysPath).
		AddCrowerExec("list --syspath", nil).
		Build()
}
