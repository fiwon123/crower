package notes

import (
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// DeleteCommandNote renders the history note for a deleted command.
//
// The note carries the whole command, not only its name: a delete is replayed
// by re-creating it, and that needs the aliases and the line the command ran.
func DeleteCommandNote(command catalog.Command, args []string) string {
	return note.New().
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.Aliases).
		AddCommandExec(command.Exec).
		AddCrowerExec("delete", args).
		Build()
}

// DeleteLastNote renders the history note for `delete --create`, `--update` or
// `--execute`.
//
// last names the operation whose command was deleted; an operation without such
// a form records the command and nothing to replay.
func DeleteLastNote(op operation.Operation, command catalog.Command) string {
	builder := note.New().
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.Aliases).
		AddCommandExec(command.Exec)

	switch op {
	case operation.Create:
		builder.AddCrowerExec("delete --create", nil)
	case operation.Update:
		builder.AddCrowerExec("delete --update", nil)
	case operation.Execute:
		builder.AddCrowerExec("delete --execute", nil)
	}

	return builder.Build()
}

// DeleteFileNote renders the history note for a deleted file.
func DeleteFileNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.File).
		AddCrowerExec("delete", args).
		Build()
}

// DeleteFolderNote renders the history note for a deleted folder.
func DeleteFolderNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.Folder).
		AddCrowerExec("delete", args).
		Build()
}

// DeleteSystemVariableNote renders the history note for `delete --system`.
func DeleteSystemVariableNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.System).
		AddCrowerExec("delete --system", args).
		Build()
}

// DeleteSystemPathNote renders the history note for `delete --syspath`.
func DeleteSystemPathNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Delete).
		AddSubOperation(operation.SysPath).
		AddCrowerExec("delete --syspath", args).
		Build()
}
