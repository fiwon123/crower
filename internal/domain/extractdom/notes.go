package extractdom

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/shell"
)

// Create a new extract note
func newExtractNote(out string, args []string) string {
	noteBuilder := note.New()

	noteBuilder.AddMainOperation(operations.Extract)

	if out != "" {
		// The --out value is quoted for the recording OS like any other value;
		// a hand-written double quote is not safe on either platform.
		noteBuilder.AddCrowerExec("extract --out "+shell.Quote(out), args)
	} else {
		noteBuilder.AddCrowerExec("extract", args)
	}

	return noteBuilder.Build()
}
