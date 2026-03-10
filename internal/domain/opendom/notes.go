package opendom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new open note
func newOpenNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Open).
		AddCrowerExec("open", args).
		Build()
}

func newOpenFolderNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Open).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("open --folder", args).
		Build()
}

func newOpenSystemNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Open).
		AddSubOperation(operationsdata.System).
		AddCrowerExec("open --system", nil).
		Build()
}
