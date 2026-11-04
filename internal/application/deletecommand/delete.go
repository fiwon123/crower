// Package deletecommand removes a managed command, a file, a folder, an
// environment variable, a PATH entry or a history entry.
//
// `delete` looks like one operation with eight flags; underneath it is five
// different jobs that happen to share a verb, a history operation and a table
// of notes. They share this package because they share the CLI, but each mode
// has its own ports:
//
//   - a managed command is removed from the catalog and committed;
//   - a file or folder is handed to the platform's remover;
//   - an environment variable is handed to the platform's environment editor;
//   - a history entry is dropped from the operation log.
//
// The command layer decides which mode a run is in -- it owns the flags -- and
// this package owns what that mode does, what it asks, and what it records.
package deletecommand

import (
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/errors"
)

// Mode is which kind of thing a request deletes.
type Mode int

const (
	// ModeCommand deletes a managed command, by name or alias. An empty Key
	// asks for one.
	ModeCommand Mode = iota

	// ModeLast deletes the command the last operation of Last touched. The
	// --create/--update/--execute forms take no key and ask nothing.
	ModeLast

	// ModeFile deletes one file path.
	ModeFile

	// ModeFolder deletes one folder and everything in it.
	ModeFolder

	// ModeSystemVariable deletes one environment variable.
	ModeSystemVariable

	// ModeSystemPath deletes one entry from the user's PATH.
	ModeSystemPath

	// ModeHistory drops one entry from the operation log.
	ModeHistory
)

// Request is one delete.
type Request struct {
	// Mode is which job to run.
	Mode Mode

	// Key is the command to delete, by name or alias. Empty means "ask".
	Key string

	// Last names the operation whose command is deleted by a --last form.
	Last operation.Operation

	// Args are the positional values the command line carried. The path and
	// variable modes act on Args[0]; every mode records them in the note.
	Args []string
}

// State is the command set, the history it is recorded in, and the way to
// persist both.
type State interface {
	// Commands returns the current command set.
	Commands() catalog.Catalog

	// LastCommand returns the command name the last operation of this kind
	// touched.
	LastCommand(op operation.Operation) (string, bool)

	// Commit replaces the command set and records the operation with the note
	// that explains it.
	Commit(commands catalog.Catalog, op operation.Operation, commandName, note string) error

	// Record appends a history entry and persists the index.
	//
	// It is separate from Commit because a file, folder or variable delete
	// changes no managed command: rewriting the command file would record a
	// change that did not happen.
	Record(op operation.Operation, commandName, note string) error
}

// Output is the diagnostics stream.
type Output interface {
	Info(message string, data ...any)
	Table(commands catalog.Catalog) string
}

// Prompter asks the questions this use case has.
type Prompter interface {
	// ChooseCommand asks which row of the printed command table to delete.
	ChooseCommand() (string, error)

	// Confirm asks a yes/no question.
	Confirm(question string) (bool, error)
}

// Remover deletes files and folders with the host's own tools.
type Remover interface {
	File(path string) error
	Folder(path string) error
}

// Environment edits the variables the user's shell reads at start.
type Environment interface {
	// Variable removes one variable and reports what it did.
	Variable(name string) (string, error)

	// VariablePath removes one entry from PATH and reports what it did.
	VariablePath(path string) (string, error)
}

// Entry is one history row, as the table shows it.
type Entry struct {
	// Version is the history version the entry records.
	Version int

	// File is the snapshot the entry points at, and the key RemoveHistoryEntry
	// takes.
	File string

	// Note is the entry's note, already decoded for display: the escaping is
	// a storage detail, so the port does not expose it.
	Note string
}

// History is the operation log and the way to drop one of its entries.
type History interface {
	// HistoryTable renders the log, newest first.
	HistoryTable() string

	// ChooseHistoryEntry asks which row to drop.
	ChooseHistoryEntry() (Entry, error)

	// RemoveHistoryEntry drops the entry and the snapshot it points at, and
	// reports what it did. Only the index is rewritten: nothing about the
	// current command set changed.
	RemoveHistoryEntry(entry Entry) (string, error)
}

// Notes renders the history note for each mode.
type Notes interface {
	DeleteCommand(command catalog.Command, args []string) string
	DeleteLast(op operation.Operation, command catalog.Command) string
	DeleteFile(args []string) string
	DeleteFolder(args []string) string
	DeleteSystemVariable(args []string) string
	DeleteSystemPath(args []string) string
}

// UseCase runs one delete.
type UseCase struct {
	state       State
	output      Output
	prompt      Prompter
	remover     Remover
	environment Environment
	history     History
	notes       Notes
}

// New builds the use case from its seven dependencies.
func New(state State, output Output, prompt Prompter, remover Remover, environment Environment, history History, notes Notes) *UseCase {
	return &UseCase{
		state:       state,
		output:      output,
		prompt:      prompt,
		remover:     remover,
		environment: environment,
		history:     history,
		notes:       notes,
	}
}

// Run deletes what the request names.
func (u *UseCase) Run(request Request) error {
	switch request.Mode {
	case ModeLast:
		return u.deleteLast(request.Last)

	case ModeFile:
		return u.deleteFile(request)

	case ModeFolder:
		return u.deleteFolder(request)

	case ModeSystemVariable:
		return u.deleteSystemVariable(request)

	case ModeSystemPath:
		return u.deleteSystemPath(request)

	case ModeHistory:
		return u.deleteHistory()

	default:
		return u.deleteCommand(request)
	}
}

// deleteCommand removes one managed command, asking which one when the request
// names none.
func (u *UseCase) deleteCommand(request Request) error {
	key := request.Key

	if key == "" {
		chosen, err := u.chooseCommand()
		if err != nil {
			return err
		}

		key = chosen
	}

	command, found := u.state.Commands().Resolve(key)
	if !found {
		u.output.Info("Command not found, try to select one.")

		chosen, err := u.chooseCommand()
		if err != nil {
			return err
		}

		key = chosen

		if command, found = u.state.Commands().Resolve(key); !found {
			return errors.CommandNotFound(key)
		}
	}

	u.showCommand(command)

	proceed, err := u.prompt.Confirm("Continue to delete")
	if err != nil {
		return err
	}

	if !proceed {
		u.output.Info("Cancelling delete...")
		return nil
	}

	return u.remove(key, func(deleted catalog.Command) string {
		return u.notes.DeleteCommand(deleted, request.Args)
	})
}

// deleteLast removes the command the last operation of a kind touched.
//
// The --last forms take no key and ask nothing: the flags describe the whole
// request, and a prompt in the middle of a scripted call would hang it.
func (u *UseCase) deleteLast(op operation.Operation) error {
	key, found := u.state.LastCommand(op)
	if !found {
		return errors.CommandNotFound(string(op))
	}

	return u.remove(key, func(deleted catalog.Command) string {
		return u.notes.DeleteLast(op, deleted)
	})
}

// remove takes the command out of the catalog, commits it with the note the
// caller's mode renders, and reports it.
func (u *UseCase) remove(key string, note func(deleted catalog.Command) string) error {
	commands, deleted, found := u.state.Commands().Remove(key)
	if !found {
		return errors.CommandNotFound(key)
	}

	if err := u.state.Commit(commands, operation.Delete, deleted.Name, note(deleted)); err != nil {
		return err
	}

	u.output.Info("deleted command: ", key)

	return nil
}

// chooseCommand prints the table and asks for a row.
func (u *UseCase) chooseCommand() (string, error) {
	u.output.Info(u.output.Table(u.state.Commands()))

	return u.prompt.ChooseCommand()
}

// showCommand prints the fields the confirmation is about.
func (u *UseCase) showCommand(command catalog.Command) {
	u.output.Info("-----------------------------------------")
	u.output.Info("Name:    ", command.Name)
	u.output.Info("Aliases: ", command.Aliases)
	u.output.Info("Exec:    ", command.Exec)
	u.output.Info("")
}

// deleteFile removes one file path and records the operation.
func (u *UseCase) deleteFile(request Request) error {
	if len(request.Args) == 0 {
		return errors.NotFileAndOutputPath()
	}

	if err := u.remover.File(request.Args[0]); err != nil {
		return err
	}

	return u.state.Record(operation.Delete, "", u.notes.DeleteFile(request.Args))
}

// deleteFolder removes one folder and records the operation.
func (u *UseCase) deleteFolder(request Request) error {
	if len(request.Args) == 0 {
		return errors.NotFileAndOutputPath()
	}

	if err := u.remover.Folder(request.Args[0]); err != nil {
		return err
	}

	return u.state.Record(operation.Delete, "", u.notes.DeleteFolder(request.Args))
}

// deleteSystemVariable removes one environment variable and records it.
func (u *UseCase) deleteSystemVariable(request Request) error {
	if len(request.Args) == 0 {
		return errors.NotEnoughArgs("var name")
	}

	out, err := u.environment.Variable(request.Args[0])
	if err != nil {
		return err
	}

	u.output.Info(out)

	return u.state.Record(operation.Delete, "", u.notes.DeleteSystemVariable(request.Args))
}

// deleteSystemPath removes one PATH entry and records it.
func (u *UseCase) deleteSystemPath(request Request) error {
	if len(request.Args) == 0 {
		return errors.NotEnoughArgs("path")
	}

	out, err := u.environment.VariablePath(request.Args[0])
	if err != nil {
		return err
	}

	u.output.Info(out)

	return u.state.Record(operation.Delete, "", u.notes.DeleteSystemPath(request.Args))
}

// deleteHistory drops one entry from the log, asking which one.
func (u *UseCase) deleteHistory() error {
	u.output.Info(u.history.HistoryTable())

	entry, err := u.history.ChooseHistoryEntry()
	if err != nil {
		return err
	}

	u.output.Info("-----------------------------------------")
	u.output.Info("Content")
	u.output.Info("Version:    ", entry.Version)
	u.output.Info("File:    ", entry.File)
	u.output.Info("Note:    ", entry.Note)
	u.output.Info("")

	proceed, err := u.prompt.Confirm("Continue to delete history content")
	if err != nil {
		return err
	}

	if !proceed {
		u.output.Info("Cancelling Delete History Content...")
		return nil
	}

	out, err := u.history.RemoveHistoryEntry(entry)
	if err != nil {
		return err
	}

	u.output.Info(out)

	return nil
}
