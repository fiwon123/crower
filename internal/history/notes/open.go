package notes

import "github.com/fiwon123/crower/internal/data/state"

// Create a new open note
func GenerateOpenNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Open).
		AddCrowerExec("open", args).
		Build()
}

func GenerateOpenFolderNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Open).
		AddSubOperation(state.Folder).
		AddCrowerExec("open --folder", args).
		Build()
}

func GenerateOpenSystemNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Open).
		AddSubOperation(state.System).
		AddCrowerExec("open --system", nil).
		Build()
}
