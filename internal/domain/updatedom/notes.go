package updatedom

import (
	"encoding/csv"
	"slices"
	"strings"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"

	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/shell"
)

// updateFlags records the values the update left behind, not only the changed
// ones: the update prompt asks for every field that was not supplied, so a
// line carrying just the delta would still prompt for the rest. Aliases are
// omitted when there are none on either side (an empty --alias parses as "not
// supplied" anyway); an update that clears aliases still cannot replay without
// a prompt, which is recorded in the issue.
func updateFlags(oldCommand *command.Data, newCommand *command.Data) string {
	var flags strings.Builder

	flags.WriteString(" --name ")
	flags.WriteString(shell.Quote(newCommand.Name))

	if len(newCommand.AllAlias) > 0 || len(oldCommand.AllAlias) > 0 {
		flags.WriteString(" --alias ")
		flags.WriteString(shell.Quote(csvAliases(newCommand.AllAlias)))
	}

	flags.WriteString(" --exec ")
	flags.WriteString(shell.Quote(newCommand.Exec))

	return flags.String()
}

// csvAliases encodes the list the way the StringSlice flag parses it, so one
// alias containing a comma survives replay.
func csvAliases(aliases []string) string {
	var list strings.Builder
	writer := csv.NewWriter(&list)

	if err := writer.Write(aliases); err != nil {
		return ""
	}

	writer.Flush()

	return strings.TrimSuffix(list.String(), "\n")
}

// Create a new update note. The positional key is the command's final name: an
// alias the update dropped (or a name it renamed) would not resolve on replay.
func newUpdateCommmandNote(oldCommand *command.Data, newCommand *command.Data) string {

	noteBuilder := note.New()

	noteBuilder.
		AddMainOperation(operations.Update).
		AddSubOperation(operations.Command)

	addCommandDiff(noteBuilder, oldCommand, newCommand)

	return noteBuilder.
		AddCrowerExec("update"+updateFlags(oldCommand, newCommand), []string{newCommand.Name}).
		Build()
}

func newUpdateLastNote(op operations.MainOperationEnum, oldCommand *command.Data, newCommand *command.Data) string {
	noteBuilder := note.New()

	noteBuilder.
		AddMainOperation(operations.Update).
		AddSubOperation(operations.Command)

	addCommandDiff(noteBuilder, oldCommand, newCommand)

	// The --last forms take no positional key; the flags still apply to the
	// command the last operation touched.
	switch op {
	case operations.Update:
		noteBuilder.AddCrowerExec("update --last"+updateFlags(oldCommand, newCommand), nil)
	case operations.Create:
		noteBuilder.AddCrowerExec("update --create"+updateFlags(oldCommand, newCommand), nil)
	case operations.Execute:
		noteBuilder.AddCrowerExec("update --execute"+updateFlags(oldCommand, newCommand), nil)
	}

	return noteBuilder.Build()
}

func addCommandDiff(noteBuilder *note.Note, oldCommand *command.Data, newCommand *command.Data) {
	if oldCommand.Name != newCommand.Name {
		noteBuilder.
			AddOldCommandName(oldCommand.Name).
			AddCommandName(newCommand.Name)
	}

	if hasDiffAlieses(oldCommand.AllAlias, newCommand.AllAlias) {
		noteBuilder.
			AddOldCommandAlias(oldCommand.AllAlias).
			AddCommandAlias(newCommand.AllAlias)
	}

	if oldCommand.Exec != newCommand.Exec {
		noteBuilder.
			AddOldCommandExec(oldCommand.Exec).
			AddCommandExec(newCommand.Exec)
	}
}

func hasDiffAlieses(old []string, new []string) bool {
	oldSorted := slices.Clone(old)
	newSorted := slices.Clone(new)

	slices.Sort(oldSorted)
	slices.Sort(newSorted)

	return !slices.Equal(oldSorted, newSorted)
}
