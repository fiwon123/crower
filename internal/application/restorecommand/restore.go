// Package restorecommand brings a command back from a history snapshot.
//
// The user picks an entry, the key they typed is resolved inside that snapshot,
// and the command is updated when its name exists and created otherwise. The
// snapshot reader and the command store are behind the State port.
package restorecommand

import (
	"fmt"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Request is one restore.
type Request struct {
	// Key is the command to look for in the snapshot, by name or alias. Empty
	// matches nothing, which is what an interactive run with no argument does.
	Key string

	// Args are the positional values, recorded in the note.
	Args []string
}

// Entry is one history row, as the table shows it.
type Entry struct {
	Version int
	File    string
	Note    string
}

// State is the command set, the snapshot lookup, and the two writes.
type State interface {
	// Commands returns the current command set.
	Commands() catalog.Catalog

	// HistoryRows returns the log newest first.
	HistoryRows() []Entry

	// FindCommand resolves key inside the snapshot file.
	FindCommand(key, file string) (catalog.Command, error)

	// Save replaces the stored command set.
	Save(commands catalog.Catalog) error

	// Record appends a history entry and persists the index.
	Record(op operation.Operation, commandName, note string) error
}

// Output is the diagnostics stream and the history table.
type Output interface {
	Info(message string, data ...any)
	Warning(message string, data ...any)
	HistoryTable() string
}

// Prompter asks the questions this use case has.
type Prompter interface {
	ChooseHistoryEntry() (Entry, error)
	Confirm(question string) (bool, error)
}

// Notes renders the history note for the restore.
type Notes interface {
	RestoreNote(args []string) string
}

// UseCase runs one restore.
type UseCase struct {
	state  State
	output Output
	prompt Prompter
	notes  Notes
}

// New builds the use case from its four dependencies.
func New(state State, output Output, prompt Prompter, notes Notes) *UseCase {
	return &UseCase{state: state, output: output, prompt: prompt, notes: notes}
}

// Run restores the command the chosen snapshot holds.
func (u *UseCase) Run(request Request) error {
	entry, proceed, err := u.choose()
	if err != nil {
		return err
	}

	if !proceed {
		u.output.Info("Cancelling restore...")
		return nil
	}

	restored, err := u.state.FindCommand(request.Key, entry.File)
	if err != nil {
		return err
	}

	commands := u.state.Commands()

	if _, exists := commands.Get(restored.Name); exists {
		return u.update(commands, request, restored)
	}

	next, err := commands.Add(restored)
	if err != nil {
		return err
	}

	if err := u.state.Save(next); err != nil {
		return err
	}

	u.output.Info("restored command: ", "restored by creating: "+describe(restored))

	return u.state.Record(operation.Restore, "", u.notes.RestoreNote(request.Args))
}

// update rewrites the existing command's exec from the snapshot.
func (u *UseCase) update(commands catalog.Catalog, request Request, restored catalog.Command) error {
	result, err := commands.Update(request.Key, catalog.Changes{Exec: restored.Exec})
	if err != nil {
		return err
	}

	for _, warning := range result.Warnings {
		u.output.Warning(warning)
	}

	if err := u.state.Save(result.Catalog); err != nil {
		return err
	}

	u.output.Info("restored command: ", fmt.Sprintf("restored by updating: from %s to %s", describe(result.Old), describe(result.New)))

	return u.state.Record(operation.Restore, "", u.notes.RestoreNote(request.Args))
}

// choose asks which history entry to restore and for confirmation.
func (u *UseCase) choose() (Entry, bool, error) {
	u.output.Info(u.output.HistoryTable())

	entry, err := u.prompt.ChooseHistoryEntry()
	if err != nil {
		return Entry{}, false, err
	}

	u.output.Info("-----------------------------------------")
	u.output.Info("Content")
	u.output.Info("Version:    ", entry.Version)
	u.output.Info("File:    ", entry.File)
	u.output.Info("Note:    ", entry.Note)
	u.output.Info("")

	proceed, err := u.prompt.Confirm("Continue to restore")
	if err != nil {
		return Entry{}, false, err
	}

	return entry, proceed, nil
}

// describe renders a command the way the restore message always has: the
// stored shape, values separated by spaces.
func describe(command catalog.Command) string {
	return fmt.Sprintf("&{%s %v %s %v}", command.Name, command.Aliases, command.Exec, command.Steps)
}
