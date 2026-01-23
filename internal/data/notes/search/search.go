package searchnotesdata

import (
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new search browser note
func NewSearchBrowserNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Search).
		AddSubOperation(operationsdata.Browser).
		AddCrowerExec("search --browser", args).
		Build()
}

func NewSearchFileNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Search).
		AddSubOperation(operationsdata.File).
		AddCrowerExec("search --file", args).
		Build()
}

func NewSearchFolderNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Search).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("search --folder", args).
		Build()
}

func NewSearchFileAndFolderNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Search).
		AddSubOperation(operationsdata.FileAndFolder).
		AddCrowerExec("search", args).
		Build()
}
