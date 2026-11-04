package notes

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// SearchBrowserNote renders the history note for `search --browser`.
func SearchBrowserNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Search).
		AddSubOperation(operation.Browser).
		AddCrowerExec("search --browser", args).
		Build()
}

// SearchFileNote renders the history note for `search --file`.
func SearchFileNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Search).
		AddSubOperation(operation.File).
		AddCrowerExec("search --file", args).
		Build()
}

// SearchFolderNote renders the history note for `search --folder`.
func SearchFolderNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Search).
		AddSubOperation(operation.Folder).
		AddCrowerExec("search --folder", args).
		Build()
}

// SearchBothNote renders the history note for the positional search.
func SearchBothNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Search).
		AddSubOperation(operation.FileAndFolder).
		AddCrowerExec("search", args).
		Build()
}
