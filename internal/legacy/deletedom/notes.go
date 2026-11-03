package deletedom

import (
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"

	"github.com/fiwon123/crower/internal/domain/operation"
)

// Create a new delete command note
func newDeleteCommandNote(command *command.Data, args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("delete", args).
		Build()
}

func newDeleteLastNote(op operation.Operation, command *command.Data) string {
	noteBuilder := note.New()

	noteBuilder.
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec)

	switch op {
	case operation.Create:
		noteBuilder.AddCrowerExec("delete --create", nil)
	case operation.Update:
		noteBuilder.AddCrowerExec("delete --update", nil)
	case operation.Execute:
		noteBuilder.AddCrowerExec("delete --execute", nil)
	}

	return noteBuilder.Build()
}

func newDeleteFileNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.File).
		AddCrowerExec("delete", args).
		Build()
}

func generateDeleteFolderNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.Folder).
		AddCrowerExec("delete", args).
		Build()
}

func newDeleteSystemVariable(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.System).
		AddCrowerExec("delete --system", args).
		Build()
}

func newDeleteSystemPathVariable(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.SysPath).
		AddCrowerExec("delete --syspath", args).
		Build()
}
