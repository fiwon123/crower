package listnotesdata

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new list note
func NewListCommandsNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.Command).
		AddCrowerExec("list", nil).
		Build()
}

func NewListProcessNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.Process).
		AddCrowerExec("list --process", args).
		Build()
}

func NewListHistoriesNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.History).
		AddCrowerExec("list --history", nil).
		Build()
}

func NewListFolderNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("list --folder", args).
		Build()
}

func NewListSystemNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.System).
		AddCrowerExec("list --system", nil).
		Build()
}

func NewListSystemPathNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.SysPath).
		AddCrowerExec("list --syspath", nil).
		Build()
}
