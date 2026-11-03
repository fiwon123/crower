package searchdom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new search browser note
func newSearchBrowserNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Search).
		AddSubOperation(operation.Browser).
		AddCrowerExec("search --browser", args).
		Build()
}

func newSearchFileNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Search).
		AddSubOperation(operation.File).
		AddCrowerExec("search --file", args).
		Build()
}

func newSearchFolderNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Search).
		AddSubOperation(operation.Folder).
		AddCrowerExec("search --folder", args).
		Build()
}

func newSearchFileAndFolderNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Search).
		AddSubOperation(operation.FileAndFolder).
		AddCrowerExec("search", args).
		Build()
}
