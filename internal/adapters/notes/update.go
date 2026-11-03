// Package notes renders the note a history entry carries.
//
// A note is both an audit line and a replay recipe: `crower_exec=` holds the
// command line that reproduces the operation, and the paired fields hold what
// changed. The format is user-visible (it is stored in historyhelper.json and
// printed by `crower list --history`) and parsers exist for it, so it lives
// here, on the storage side, rather than in a command.
package notes

import (
	"encoding/csv"
	"slices"
	"strings"

	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/shell"
)

// UpdateNote renders the history note for an update.
//
// last names the operation whose command was updated by a `--last` form; an
// empty last means the positional form, which is the only one that carries the
// command name.
func UpdateNote(before, after catalog.Command, last operation.Operation) string {
	builder := note.New()

	builder.
		AddMainOperation(operation.Update).
		AddSubOperation(operation.Command)

	addDiff(builder, before, after)

	flags := updateFlags(before, after)

	if last == "" {
		// The positional key is the command's final name: an alias the update
		// dropped, or a name it renamed, would not resolve on replay.
		return builder.
			AddCrowerExec("update"+flags, []string{after.Name}).
			Build()
	}

	// The --last forms take no positional key; the flags still apply to the
	// command the last operation touched.
	switch last {
	case operation.Update:
		return builder.AddCrowerExec("update --last"+flags, nil).Build()
	case operation.Create:
		return builder.AddCrowerExec("update --create"+flags, nil).Build()
	case operation.Execute:
		return builder.AddCrowerExec("update --execute"+flags, nil).Build()
	}

	// An operation without a --last form records the diff and nothing to
	// replay, which is what it did before this package existed.
	return builder.Build()
}

// updateFlags records the values the update left behind, not only the changed
// ones: the update prompt asks for every field that was not supplied, so a line
// carrying just the delta would still prompt for the rest. Aliases are omitted
// when there are none on either side (an empty --alias parses as "not supplied"
// anyway); an update that clears aliases still cannot replay without a prompt.
func updateFlags(before, after catalog.Command) string {
	var flags strings.Builder

	flags.WriteString(" --name ")
	flags.WriteString(shell.Quote(after.Name))

	if len(after.Aliases) > 0 || len(before.Aliases) > 0 {
		flags.WriteString(" --alias ")
		flags.WriteString(shell.Quote(csvAliases(after.Aliases)))
	}

	flags.WriteString(" --exec ")
	flags.WriteString(shell.Quote(after.Exec))

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

func addDiff(builder *note.Note, before, after catalog.Command) {
	if before.Name != after.Name {
		builder.
			AddOldCommandName(before.Name).
			AddCommandName(after.Name)
	}

	if differentAliases(before.Aliases, after.Aliases) {
		builder.
			AddOldCommandAlias(before.Aliases).
			AddCommandAlias(after.Aliases)
	}

	if before.Exec != after.Exec {
		builder.
			AddOldCommandExec(before.Exec).
			AddCommandExec(after.Exec)
	}
}

// differentAliases compares the alias lists as sets, without sorting the
// slices the command holds.
func differentAliases(before, after []string) bool {
	beforeSorted := slices.Clone(before)
	afterSorted := slices.Clone(after)

	slices.Sort(beforeSorted)
	slices.Sort(afterSorted)

	return !slices.Equal(beforeSorted, afterSorted)
}
