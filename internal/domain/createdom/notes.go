package createdom

import (
	commanddata "github.com/fiwon123/crower/internal/data/command"
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new add note
func newCreateCommandNote(command *commanddata.Data, args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func NewCreateScriptCommandNote(command *commanddata.Data, args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.Script).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func NewCreateCommandLastExecuteNote(command *commanddata.Data) string {
	noteBuilder := notesdata.New()

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
func NewCreateProcessNote(command *commanddata.Data, args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.Process).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func NewCreateSystemVariableNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.System).
		AddCrowerExec("create --system", args).
		Build()
}

func NewCreateSystemPathVariableNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.SysPath).
		AddCrowerExec("create --syspath", args).
		Build()
}

func GenerateCreateFile(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.File).
		AddCrowerExec("create", args).
		Build()
}

func NewCreateFolder(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operations.Create).
		AddSubOperation(operations.Folder).
		AddCrowerExec("create", args).
		Build()
}
