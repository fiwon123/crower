package notes

import (
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/state"
)

// Create a new delete command note
func GenerateDeleteCommandNote(command *command.Data, args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Delete).
		AddSubOperation(state.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("delete", args).
		Build()
}

func GenerateDeleteLastNote(op state.MainOperationEnum, command *command.Data) string {
	noteBuilder := New()

	noteBuilder.
		AddMainOperation(state.Delete).
		AddSubOperation(state.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec)

	switch op {
	case state.Create:
		noteBuilder.AddCrowerExec("delete --create", nil)
	case state.Update:
		noteBuilder.AddCrowerExec("delete --update", nil)
	case state.Execute:
		noteBuilder.AddCrowerExec("delete --execute", nil)
	}

	return noteBuilder.Build()
}

func GenerateDeleteFileNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Delete).
		AddSubOperation(state.File).
		AddCrowerExec("delete", args).
		Build()
}

func GenerateDeleteFolderNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Delete).
		AddSubOperation(state.Folder).
		AddCrowerExec("delete", args).
		Build()
}

func GenerateDeleteSystemVariable(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Delete).
		AddSubOperation(state.System).
		AddCrowerExec("delete --system", args).
		Build()
}

func GenerateDeleteSystemPathVariable(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Delete).
		AddSubOperation(state.SysPath).
		AddCrowerExec("delete --syspath", args).
		Build()
}
