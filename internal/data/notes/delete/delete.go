package deletenotesdata

import (
	commanddata "github.com/fiwon123/crower/internal/data/command"
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Create a new delete command note
func NewDeleteCommandNote(command *commanddata.Data, args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Delete).
		AddSubOperation(operationsdata.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("delete", args).
		Build()
}

func NewDeleteLastNote(op operationsdata.MainOperationEnum, command *commanddata.Data) string {
	noteBuilder := notesdata.New()

	noteBuilder.
		AddMainOperation(operationsdata.Delete).
		AddSubOperation(operationsdata.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec)

	switch op {
	case operationsdata.Create:
		noteBuilder.AddCrowerExec("delete --create", nil)
	case operationsdata.Update:
		noteBuilder.AddCrowerExec("delete --update", nil)
	case operationsdata.Execute:
		noteBuilder.AddCrowerExec("delete --execute", nil)
	}

	return noteBuilder.Build()
}

func NewDeleteFileNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Delete).
		AddSubOperation(operationsdata.File).
		AddCrowerExec("delete", args).
		Build()
}

func GenerateDeleteFolderNote(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Delete).
		AddSubOperation(operationsdata.Folder).
		AddCrowerExec("delete", args).
		Build()
}

func NewDeleteSystemVariable(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Delete).
		AddSubOperation(operationsdata.System).
		AddCrowerExec("delete --system", args).
		Build()
}

func NewDeleteSystemPathVariable(args []string) string {
	noteBuilder := notesdata.New()

	return noteBuilder.
		AddMainOperation(operationsdata.Delete).
		AddSubOperation(operationsdata.SysPath).
		AddCrowerExec("delete --syspath", args).
		Build()
}
