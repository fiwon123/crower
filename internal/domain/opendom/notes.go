package opendom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new open note
func newOpenNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Open).
		AddCrowerExec("open", args).
		Build()
}

func newOpenFolderNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Open).
		AddSubOperation(operations.Folder).
		AddCrowerExec("open --folder", args).
		Build()
}

func newOpenSystemNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Open).
		AddSubOperation(operations.System).
		AddCrowerExec("open --system", nil).
		Build()
}
