package notes

import (
	"fmt"

	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new extract note
func GenerateExtractNote(out string, args []string) string {
	noteBuilder := New()

	noteBuilder.AddMainOperation(operationsdata.Extract)

	if out != "" {
		noteBuilder.AddCrowerExec(fmt.Sprintf("extract --out \"%s\"", out), args)
	} else {
		noteBuilder.AddCrowerExec("extract", args)
	}

	return noteBuilder.Build()
}
