package notes

import (
	commanddata "github.com/fiwon123/crower/internal/data/command"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new add note
func GenerateCreateCommandNote(command *commanddata.Data, args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func GenerateCreateScriptCommandNote(command *commanddata.Data, args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.Script).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func GenerateCreateCommandLastExecuteNote(command *commanddata.Data) string {
	noteBuilder := New()

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
func GenerateCreateProcessNote(command *commanddata.Data, args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.Process).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create", args).
		Build()
}

func GenerateCreateSystemVariableNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.System).
		AddCrowerExec("create --system", args).
		Build()
}

func GenerateCreateSystemPathVariableNote(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.SysPath).
		AddCrowerExec("create --syspath", args).
		Build()
}

func GenerateCreateFile(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.File).
		AddCrowerExec("create", args).
		Build()
}

func GenerateCreateFolder(args []string) string {
	noteBuilder := New()

	return noteBuilder.
		AddMainOperation(operationsdata.Create).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("create", args).
		Build()
}
