package notes

import "github.com/fiwon123/crower/internal/data/state"

// Create a new search browser note
func GenerateSearchBrowserNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Search).
		AddSubOperation(state.Browser).
		AddCrowerExec("search --browser", args).
		Build()
}

func GenerateSearchFileNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Search).
		AddSubOperation(state.File).
		AddCrowerExec("search --file", args).
		Build()
}

func GenerateSearchFolderNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Search).
		AddSubOperation(state.Folder).
		AddCrowerExec("search --folder", args).
		Build()
}

func GenerateSearchFileAndFolderNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Search).
		AddSubOperation(state.FileAndFolder).
		AddCrowerExec("search", args).
		Build()
}
