package notes

import operationsdata "github.com/fiwon123/crower/internal/data/operations"

// Create a new search browser note
func GenerateSearchBrowserNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Search).
		AddSubOperation(operationsdata.Browser).
		AddCrowerExec("search --browser", args).
		Build()
}

func GenerateSearchFileNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Search).
		AddSubOperation(operationsdata.File).
		AddCrowerExec("search --file", args).
		Build()
}

func GenerateSearchFolderNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Search).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("search --folder", args).
		Build()
}

func GenerateSearchFileAndFolderNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Search).
		AddSubOperation(operationsdata.FileAndFolder).
		AddCrowerExec("search", args).
		Build()
}
