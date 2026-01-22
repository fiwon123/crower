package opennotesdata

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new open note
func NewOpenNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Open).
		AddCrowerExec("open", args).
		Build()
}

func NewOpenFolderNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Open).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("open --folder", args).
		Build()
}

func NewOpenSystemNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Open).
		AddSubOperation(operationsdata.System).
		AddCrowerExec("open --system", nil).
		Build()
}
