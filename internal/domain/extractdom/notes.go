package extractdom

import (
	"fmt"

	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new extract note
func newExtractNote(out string, args []string) string {
	noteBuilder := notesdata.New()

	noteBuilder.AddMainOperation(operations.Extract)

	if out != "" {
		noteBuilder.AddCrowerExec(fmt.Sprintf("extract --out \"%s\"", out), args)
	} else {
		noteBuilder.AddCrowerExec("extract", args)
	}

	return noteBuilder.Build()
}
