// Package reverthistory moves the command set back through the history.
package reverthistory

import (
	"fmt"
	"strconv"

	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/errors"
)

// History is the timeline and the state it restores.
type History interface {
	// Pending renders the entries a revert would drop and the ones that would
	// remain, and fails when the step count is out of range.
	Pending(steps int) (dropped, kept string, err error)

	// RevertTo restores the state steps entries back and returns the version it
	// landed on.
	RevertTo(steps int) (version int, err error)

	Record(op operation.Operation, commandName, note string) error
}

// Output is the diagnostics stream.
type Output interface {
	Info(message string, data ...any)
}

// Prompter asks the confirmation.
type Prompter interface {
	Confirm(question string) (bool, error)
}

// Notes renders the history note for the revert.
type Notes interface {
	Revert(args []string) string
}

// UseCase reverts the command set.
type UseCase struct {
	history History
	output  Output
	prompt  Prompter
	notes   Notes
}

// New builds the use case from its four dependencies.
func New(history History, output Output, prompt Prompter, notes Notes) *UseCase {
	return &UseCase{history: history, output: output, prompt: prompt, notes: notes}
}

// Run reverts the number of steps given.
func (useCase *UseCase) Run(args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("steps")
	}

	steps, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("steps must be a number: %q", args[0])
	}

	dropped, kept, err := useCase.history.Pending(steps)
	if err != nil {
		return err
	}

	useCase.output.Info("")
	useCase.output.Info("Deleted History")
	useCase.output.Info(dropped)
	useCase.output.Info("")
	useCase.output.Info("New History")
	useCase.output.Info(kept)
	useCase.output.Info("")
	useCase.output.Info("History will revert ", steps)

	proceed, err := useCase.prompt.Confirm("Continue to revert")
	if err != nil {
		return err
	}

	if !proceed {
		useCase.output.Info("Cancelling revert...")
		return nil
	}

	version, err := useCase.history.RevertTo(steps)
	if err != nil {
		return fmt.Errorf("revert history: %w", err)
	}

	useCase.output.Info("reverted to history version ", version)

	return useCase.history.Record(operation.Revert, "", useCase.notes.Revert(args))
}
