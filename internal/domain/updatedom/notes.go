package updatedom

import (
	"sort"

	commanddata "github.com/fiwon123/crower/internal/data/command"
	notesdata "github.com/fiwon123/crower/internal/data/notes"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Create a new update note
func newUpdateCommmandNote(args []string, oldCommand *commanddata.Data, newCommand *commanddata.Data) string {

	noteBuilder := notesdata.New()

	noteBuilder.
		AddMainOperation(operations.Update).
		AddSubOperation(operations.Command)

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

	return noteBuilder.
		AddCrowerExec("update", args).
		Build()
}

func newUpdateLastNote(op operations.MainOperationEnum, oldCommand *commanddata.Data, newCommand *commanddata.Data) string {
	noteBuilder := notesdata.New()

	noteBuilder.
		AddMainOperation(operations.Update).
		AddSubOperation(operations.Command)

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

	switch op {
	case operations.Update:
		noteBuilder.AddCrowerExec("update --last", nil)
	case operations.Create:
		noteBuilder.AddCrowerExec("update --create", nil)
	case operations.Execute:
		noteBuilder.AddCrowerExec("update --execute", nil)
	}

	return noteBuilder.Build()
}

func hasDiffAlieses(old []string, new []string) bool {
	lenOld := len(old)
	if lenOld != len(new) {
		return true
	}

	sort.Strings(old)
	sort.Strings(new)

	for i := range lenOld {
		if old[i] != new[i] {
			return true
		}
	}

	return false
}
