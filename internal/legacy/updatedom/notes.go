package updatedom

import (
	"github.com/fiwon123/crower/internal/adapters/notes"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// newUpdateCommmandNote renders the history note for a positional update.
//
// The format lives in internal/adapters/notes; this is the remaining bridge to
// the legacy command shape, and it goes with this package.
func newUpdateCommmandNote(oldCommand *command.Data, newCommand *command.Data) string {
	return notes.UpdateNote(toCatalog(oldCommand), toCatalog(newCommand), "")
}

// newUpdateLastNote renders the history note for one of the --last forms.
func newUpdateLastNote(op operation.Operation, oldCommand *command.Data, newCommand *command.Data) string {
	return notes.UpdateNote(toCatalog(oldCommand), toCatalog(newCommand), op)
}

func toCatalog(data *command.Data) catalog.Command {
	return catalog.Command{
		Name:    data.Name,
		Aliases: data.AllAlias,
		Exec:    data.Exec,
		Steps:   data.Steps,
	}
}
