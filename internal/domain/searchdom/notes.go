package searchdom

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new search browser note
func newSearchBrowserNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Search).
		AddSubOperation(operations.Browser).
		AddCrowerExec("search --browser", args).
		Build()
}

func newSearchFileNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Search).
		AddSubOperation(operations.File).
		AddCrowerExec("search --file", args).
		Build()
}

func newSearchFolderNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Search).
		AddSubOperation(operations.Folder).
		AddCrowerExec("search --folder", args).
		Build()
}

func newSearchFileAndFolderNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Search).
		AddSubOperation(operations.FileAndFolder).
		AddCrowerExec("search", args).
		Build()
}
