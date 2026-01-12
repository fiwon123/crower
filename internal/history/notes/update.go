package notes

import (
	"sort"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/state"
)

// Create a new update note
func GenerateUpdateCommmandNote(args []string, oldCommand *command.Data, newCommand *command.Data) string {

	noteBuilder := New()

	noteBuilder.
		AddMainOperation(state.Update).
		AddSubOperation(state.Command)

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

func GenerateUpdateLastNote(op state.MainOperationEnum, oldCommand *command.Data, newCommand *command.Data) string {
	noteBuilder := New()

	noteBuilder.
		AddMainOperation(state.Update).
		AddSubOperation(state.Command)

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
	case state.Update:
		noteBuilder.AddCrowerExec("update --last", nil)
	case state.Create:
		noteBuilder.AddCrowerExec("update --create", nil)
	case state.Execute:
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
