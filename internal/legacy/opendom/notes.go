package opendom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new open note
func newOpenNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Open).
		AddCrowerExec("open", args).
		Build()
}

func newOpenFolderNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Open).
		AddSubOperation(operation.Folder).
		AddCrowerExec("open --folder", args).
		Build()
}

func newOpenSystemNote() string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Open).
		AddSubOperation(operation.System).
		AddCrowerExec("open --system", nil).
		Build()
}
