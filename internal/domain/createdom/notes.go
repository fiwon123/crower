package createdom

import (
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"

	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new add note
func newCreateCommandNote(command *command.Data, args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func newCreateScriptCommandNote(command *command.Data, args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.Script).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func newCreateCommandLastExecuteNote(command *command.Data) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.Last).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create --execute", nil).
		Build()
}

// Create a new add process note
func newCreateProcessNote(command *command.Data, args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.Process).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func newCreateSystemVariableNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.System).
		AddCrowerExec("create --system", args).
		Build()
}

func newCreateSystemPathVariableNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.SysPath).
		AddCrowerExec("create --syspath", args).
		Build()
}

func generateCreateFile(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.File).
		AddCrowerExec("create", args).
		Build()
}

func newCreateFolder(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.Folder).
		AddCrowerExec("create", args).
		Build()
}
