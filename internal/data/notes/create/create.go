package createnotesdata

import (
	commanddata "github.com/fiwon123/crower/internal/data/command"
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new add note
func NewCreateCommandNote(command *commanddata.Data, args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func NewCreateScriptCommandNote(command *commanddata.Data, args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.Script).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func NewCreateCommandLastExecuteNote(command *commanddata.Data) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.Last).
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
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.Process).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func NewCreateSystemVariableNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.System).
		AddCrowerExec("create --system", args).
		Build()
}

func NewCreateSystemPathVariableNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.SysPath).
		AddCrowerExec("create --syspath", args).
		Build()
}

func GenerateCreateFile(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.File).
		AddCrowerExec("create", args).
		Build()
}

func NewCreateFolder(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("create", args).
		Build()
}
