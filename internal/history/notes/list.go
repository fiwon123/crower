package notes

import (
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new list note
func GenerateListCommandsNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.Command).
		AddCrowerExec("list", nil).
		Build()
}

func GenerateListProcessNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.Process).
		AddCrowerExec("list --process", args).
		Build()
}

func GenerateListHistoriesNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.History).
		AddCrowerExec("list --history", nil).
		Build()
}

func GenerateListFolderNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("list --folder", args).
		Build()
}

func GenerateListSystemNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.System).
		AddCrowerExec("list --system", nil).
		Build()
}

func GenerateListSystemPathNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.SysPath).
		AddCrowerExec("list --syspath", nil).
		Build()
}
