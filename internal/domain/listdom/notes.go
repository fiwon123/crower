package listdom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new list note
func newListCommandsNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.Command).
		AddCrowerExec("list", nil).
		Build()
}

func newListProcessNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.Process).
		AddCrowerExec("list --process", args).
		Build()
}

func newListHistoriesNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.History).
		AddCrowerExec("list --history", nil).
		Build()
}

func newListFolderNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("list --folder", args).
		Build()
}

func newListSystemNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.System).
		AddCrowerExec("list --system", nil).
		Build()
}

func newListSystemPathNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.List).
		AddSubOperation(operationsdata.SysPath).
		AddCrowerExec("list --syspath", nil).
		Build()
}
