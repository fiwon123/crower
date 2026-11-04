package notes

import (
	"encoding/csv"
	"strings"

	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/catalog"
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

// CreateCommandNote renders the history note for `create <name> <exec>`.
func CreateCommandNote(command catalog.Command) string {
	return note.New().
		AddMainOperation(operation.Create).
		AddSubOperation(operation.Command).
		AddCommandName(command.Name).
		AddCommandAlias(command.Aliases).
		AddCommandExec(command.Exec).
		AddCrowerExec("create"+aliasArg(command.Aliases), []string{command.Name, command.Exec}).
		Build()
}

// CreateScriptNote renders the history note for `create --script`.
func CreateScriptNote(command catalog.Command) string {
	return note.New().
		AddMainOperation(operation.Create).
		AddSubOperation(operation.Script).
		AddCommandName(command.Name).
		AddCommandAlias(command.Aliases).
		AddCommandExec(command.Exec).
		AddCrowerExec("create --script"+aliasArg(command.Aliases), []string{command.Name}).
		Build()
}

// CreateLastNote renders the history note for `create --execute`.
func CreateLastNote(command catalog.Command) string {
	return note.New().
		AddMainOperation(operation.Create).
		AddSubOperation(operation.Last).
		AddCommandName(command.Name).
		AddCommandAlias(command.Aliases).
		AddCommandExec(command.Exec).
		AddCrowerExec("create --execute"+aliasArg(command.Aliases), []string{command.Name}).
		Build()
}

// CreateProcessNote renders the history note for `create --process`.
//
// The replay line carries the new name as the flag's value and the original
// positional values after it.
func CreateProcessNote(command catalog.Command, args []string) string {
	return note.New().
		AddMainOperation(operation.Create).
		AddSubOperation(operation.Process).
		AddCommandName(command.Name).
		AddCommandAlias(command.Aliases).
		AddCommandExec(command.Exec).
		AddCrowerExec("create --process "+shell.Quote(command.Name)+aliasArg(command.Aliases), args).
		Build()
}

// CreateSystemVariableNote renders the history note for `create --system`.
func CreateSystemVariableNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Create).
		AddSubOperation(operation.System).
		AddCrowerExec("create --system", args).
		Build()
}

// CreateSystemPathNote renders the history note for `create --syspath`.
func CreateSystemPathNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Create).
		AddSubOperation(operation.SysPath).
		AddCrowerExec("create --syspath", args).
		Build()
}

// CreateFileNote renders the history note for `create --file`.
func CreateFileNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Create).
		AddSubOperation(operation.File).
		AddCrowerExec("create --file", args).
		Build()
}

// CreateFolderNote renders the history note for `create --folder`.
func CreateFolderNote(args []string) string {
	return note.New().
		AddMainOperation(operation.Create).
		AddSubOperation(operation.Folder).
		AddCrowerExec("create --folder", args).
		Build()
}
