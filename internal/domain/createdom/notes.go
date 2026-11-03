package createdom

import (
	"encoding/csv"
	"strings"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"

	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/shell"
)

// aliasArg records the --alias flag, so the line can be replayed with the
// command it created rather than one that lost its aliases.
//
// The flag is a pflag StringSlice, which reads CSV: a comma or a quote inside
// one alias has to go through csv.Writer, or replay splits it into two.
func aliasArg(aliases []string) string {
	if len(aliases) == 0 {
		return ""
	}

	var list strings.Builder
	writer := csv.NewWriter(&list)

	if err := writer.Write(aliases); err != nil {
		// csv.Writer only fails on a bad *Reader, never on Write.
		return ""
	}

	writer.Flush()

	return " --alias " + shell.Quote(strings.TrimSuffix(list.String(), "\n"))
}

// Create a new add note
func newCreateCommandNote(command *command.Data) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Create).
		AddSubOperation(operation.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create"+aliasArg(command.AllAlias), []string{command.Name, command.Exec}).
		Build()
}

func newCreateScriptCommandNote(command *command.Data) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Create).
		AddSubOperation(operation.Script).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create --script"+aliasArg(command.AllAlias), []string{command.Name}).
		Build()
}

func newCreateCommandLastExecuteNote(command *command.Data) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Create).
		AddSubOperation(operation.Last).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create --execute"+aliasArg(command.AllAlias), []string{command.Name}).
		Build()
}

// Create a new add process note
func newCreateProcessNote(command *command.Data, args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Create).
		AddSubOperation(operation.Process).
		AddCommandName(command.Name).
		AddCommandAlias(command.AllAlias).
		AddCommandExec(command.Exec).
		AddCrowerExec("create --process "+shell.Quote(command.Name)+aliasArg(command.AllAlias), args).
		Build()
}

func newCreateSystemVariableNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Create).
		AddSubOperation(operation.System).
		AddCrowerExec("create --system", args).
		Build()
}

func newCreateSystemPathVariableNote(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Create).
		AddSubOperation(operation.SysPath).
		AddCrowerExec("create --syspath", args).
		Build()
}

func generateCreateFile(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Create).
		AddSubOperation(operation.File).
		AddCrowerExec("create --file", args).
		Build()
}

func newCreateFolder(args []string) string {
	noteBuilder := note.New()

	return noteBuilder.
		AddMainOperation(operation.Create).
		AddSubOperation(operation.Folder).
		AddCrowerExec("create --folder", args).
		Build()
}
