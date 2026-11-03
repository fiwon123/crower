// Package updatecommand changes a managed command.
//
// The rules -- which name and alias changes are allowed, how an alias claim is
// repointed or kept dormant -- live in the catalog aggregate. This package owns
// the flow: which command, which fields, and what the user is asked before
// anything is written.
package updatecommand

import (
	"fmt"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/errors"
)

// Request is one update.
type Request struct {
	// Key is the command to update, by name or alias. Empty means "ask".
	Key string

	// Changes are the fields supplied on the command line. An empty field is
	// asked for, or kept.
	Changes catalog.Changes

	// Last, when set, updates the command the last operation of this kind
	// touched. The --last forms take no key and ask nothing.
	Last operation.Operation
}

// State is the command set, the history it is recorded in, and the way to
// persist both.
type State interface {
	Commands() catalog.Catalog

	// LastCommand returns the command name the last operation of this kind
	// touched.
	LastCommand(op operation.Operation) (string, bool)

	// Commit replaces the command set and records the operation with the note
	// that explains it.
	Commit(commands catalog.Catalog, op operation.Operation, commandName, note string) error
}

// Output is the diagnostics stream.
type Output interface {
	Info(message string, data ...any)
	Table(commands catalog.Catalog) string
}

// Prompter asks the user for the values an update needs.
//
// The question texts belong to the prompter: these are the prompts crower has
// always shown, and a port that restated them here would be a second place to
// change one.
type Prompter interface {
	ChooseCommand() (string, error)
	NewName() (string, error)
	NewAliases() ([]string, error)
	NewExec() (string, error)
	Confirm(question string) (bool, error)
}

// Notes renders the history note for the change.
type Notes interface {
	Update(before, after catalog.Command, last operation.Operation) string
}

// UseCase updates one command.
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

// Run applies one update.
//
// A declined confirmation is not a failure: nothing happened, and that was the
// answer.
func (u *UseCase) Run(request Request) error {
	if request.Last != "" {
		return u.updateLast(request)
	}

	key := request.Key
	changes := request.Changes

	if key == "" {
		u.output.Info(u.output.Table(u.state.Commands()))

		chosen, err := u.prompt.ChooseCommand()
		if err != nil {
			return err
		}

		key = chosen
	}

	current, found := u.state.Commands().Resolve(key)
	if !found {
		return errors.CommandNotFound(key)
	}

	u.showCurrent(current)

	changes, err := u.askForMissing(changes)
	if err != nil {
		return err
	}

	u.showChange(current, changes)

	proceed, err := u.prompt.Confirm("Continue to update")
	if err != nil {
		return err
	}

	if !proceed {
		u.output.Info("Cancelling update...")
		return nil
	}

	return u.apply(key, changes, "")
}

// updateLast applies the change to the command the last operation touched.
//
// The --last forms take no key and ask nothing: the flags describe the whole
// change, and a prompt in the middle of a scripted call would hang it.
func (u *UseCase) updateLast(request Request) error {
	key, found := u.state.LastCommand(request.Last)
	if !found {
		return fmt.Errorf("no %q command in history", request.Last)
	}

	return u.apply(key, request.Changes, request.Last)
}

// apply writes the change and records it.
func (u *UseCase) apply(key string, changes catalog.Changes, last operation.Operation) error {
	result, err := u.state.Commands().Update(key, changes)
	if err != nil {
		return err
	}

	note := u.notes.Update(result.Old, result.New, last)

	if err := u.state.Commit(result.Catalog, operation.Update, result.New.Name, note); err != nil {
		return err
	}

	u.output.Info("updated command: ", key)

	return nil
}

// askForMissing asks for every field the command line did not supply.
func (u *UseCase) askForMissing(changes catalog.Changes) (catalog.Changes, error) {
	if changes.Name == "" {
		proceed, err := u.prompt.Confirm("Do you want to update name")
		if err != nil {
			return changes, err
		}

		if proceed {
			name, err := u.prompt.NewName()
			if err != nil {
				return changes, err
			}

			changes.Name = name
		}
	}

	if len(changes.Aliases) == 0 {
		proceed, err := u.prompt.Confirm("Do you want to update alias")
		if err != nil {
			return changes, err
		}

		if proceed {
			aliases, err := u.prompt.NewAliases()
			if err != nil {
				return changes, err
			}

			changes.Aliases = aliases
		}
	}

	if changes.Exec == "" {
		proceed, err := u.prompt.Confirm("Do you want to update exec")
		if err != nil {
			return changes, err
		}

		if proceed {
			exec, err := u.prompt.NewExec()
			if err != nil {
				return changes, err
			}

			changes.Exec = exec
		}
	}

	return changes, nil
}

// showCurrent prints the command as it is, before anything is asked.
func (u *UseCase) showCurrent(command catalog.Command) {
	u.output.Info("-----------------------------------------")
	u.output.Info("Name:    ", command.Name)
	u.output.Info("Aliases: ", command.Aliases)
	u.output.Info("Exec:    ", command.Exec)
	u.output.Info("")
}

// showChange prints what the update would leave behind, for the confirmation.
func (u *UseCase) showChange(before catalog.Command, changes catalog.Changes) {
	u.output.Info("-----------------------------------------")
	u.output.Info("Old Command: ")
	u.output.Info("Name:    ", before.Name)
	u.output.Info("Aliases: ", before.Aliases)
	u.output.Info("Exec:    ", before.Exec)

	u.output.Info("-----------------------------------------")
	u.output.Info("New Command: ")
	u.output.Info("Name:    ", changes.Name)
	u.output.Info("Aliases: ", changes.Aliases)
	u.output.Info("Exec:    ", changes.Exec)
	u.output.Info("")
}
