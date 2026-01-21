package notes

import operationsdata "github.com/fiwon123/crower/internal/data/operations"

// Create a new open note
func GenerateOpenNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Open).
		AddCrowerExec("open", args).
		Build()
}

func GenerateOpenFolderNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Open).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("open --folder", args).
		Build()
}

func GenerateOpenSystemNote() string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Open).
		AddSubOperation(operationsdata.System).
		AddCrowerExec("open --system", nil).
		Build()
}
