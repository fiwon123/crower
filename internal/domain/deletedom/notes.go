package deletedom

import (
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"

	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new delete command note
func NewDeleteCommandNote(command *command.Data, args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Delete).
		AddSubOperation(operations.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("delete", args).
		Build()
}

func NewDeleteLastNote(op operations.MainOperationEnum, command *command.Data) string {
	noteBuilder := note.New()

	noteBuilder.
		AddMainOperation(operations.Delete).
		AddSubOperation(operations.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec)

	switch op {
	case operations.Create:
		noteBuilder.AddCrowerExec("delete --create", nil)
	case operations.Update:
		noteBuilder.AddCrowerExec("delete --update", nil)
	case operations.Execute:
		noteBuilder.AddCrowerExec("delete --execute", nil)
	}

	return noteBuilder.Build()
}

func newDeleteFileNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Delete).
		AddSubOperation(operations.File).
		AddCrowerExec("delete", args).
		Build()
}

func generateDeleteFolderNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Delete).
		AddSubOperation(operations.Folder).
		AddCrowerExec("delete", args).
		Build()
}

func NewDeleteSystemVariable(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Delete).
		AddSubOperation(operations.System).
		AddCrowerExec("delete --system", args).
		Build()
}

func NewDeleteSystemPathVariable(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Delete).
		AddSubOperation(operations.SysPath).
		AddCrowerExec("delete --syspath", args).
		Build()
}
