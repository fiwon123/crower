package notes

import (
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/shell"
)

// ExtractNote renders the history note for an extraction.
//
// out is the --out value, if one was given.
func ExtractNote(out string, args []string) string {
	builder := New().AddMainOperation(operation.Extract)

	if out != "" {
		// The --out value is quoted for the recording OS like any other value;
		// a hand-written double quote is not safe on either platform.
		return builder.AddCrowerExec("extract --out "+shell.Quote(out), args).Build()
	}

	return builder.AddCrowerExec("extract", args).Build()
}
