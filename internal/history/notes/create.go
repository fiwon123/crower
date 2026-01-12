package notes

import (
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/state"
)

// Create a new add note
func GenerateCreateCommandNote(command *command.Data, args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Create).
		AddSubOperation(state.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func GenerateCreateCommandLastExecuteNote(command *command.Data) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Create).
		AddSubOperation(state.Last).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create --execute", nil).
		Build()
}

// Create a new add process note
func GenerateCreateProcessNote(command *command.Data, args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Create).
		AddSubOperation(state.Process).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func GenerateCreateSystemVariableNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Create).
		AddSubOperation(state.System).
		AddCrowerExec("create --system", args).
		Build()
}

func GenerateCreateSystemPathVariableNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Create).
		AddSubOperation(state.SysPath).
		AddCrowerExec("create --syspath", args).
		Build()
}

func GenerateCreateFile(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Create).
		AddSubOperation(state.File).
		AddCrowerExec("create", args).
		Build()
}

func GenerateCreateFolder(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(state.Create).
		AddSubOperation(state.Folder).
		AddCrowerExec("create", args).
		Build()
}
