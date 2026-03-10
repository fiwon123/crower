package opendom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new open note
func newOpenNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Open).
		AddCrowerExec("open", args).
		Build()
}

func newOpenFolderNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Open).
		AddSubOperation(operations.Folder).
		AddCrowerExec("open --folder", args).
		Build()
}

func newOpenSystemNote() string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Open).
		AddSubOperation(operations.System).
		AddCrowerExec("open --system", nil).
		Build()
}
