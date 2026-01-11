package notes

import (
	"github.com/fiwon123/crower/internal/data/command"
)

// Create a new execute note
func GenerateExecuteNote(command *command.Data) string {
	noteBuilder := New()
	noteBuilder.AddCommandName(command.Name)
	noteBuilder.AddCommandAlias(command.AllAlias)
	noteBuilder.AddCommandExec(command.Exec)
	noteBuilder.AddCrowerExec("execute", nil)

	return noteBuilder.Build()
}
