package notes

import (
	"github.com/fiwon123/crower/internal/data/command"
)

// Create a new delete command note
func GenerateDeleteCommandNote(command *command.Data) string {
	noteBuilder := New()
	noteBuilder.AddCommandName(command.Name)
	noteBuilder.AddCommandAlias(command.AllAlias)
	noteBuilder.AddCommandExec(command.Exec)
	noteBuilder.AddCrowerExec("delete", nil)

	return noteBuilder.Build()
}

func GenerateDeleteFileNote(args []string) string {
	noteBuilder := New()
	noteBuilder.AddCrowerExec("delete", args)

	return noteBuilder.Build()
}

func GenerateDeleteFolderNote(args []string) string {
	noteBuilder := New()
	noteBuilder.AddCrowerExec("delete", args)

	return noteBuilder.Build()
}
