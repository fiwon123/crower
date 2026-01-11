package notes

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/state"
)

// Create a new extract note
func GenerateExtractNote(out string, args []string) string {
	noteBuilder := New()

	noteBuilder.AddMainOperation(state.Extract)

	if out != "" {
		noteBuilder.AddCrowerExec(fmt.Sprintf("extract --out \"%s\"", out), args)
	} else {
		noteBuilder.AddCrowerExec("extract", args)
	}

	return noteBuilder.Build()
}
