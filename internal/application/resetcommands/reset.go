// Package resetcommands erases every managed command.
//
// It is the first use case of the new core, and it is deliberately the smallest
// one that still needs everything a use case needs: the state, the diagnostics
// stream, and a question.
//
// The interfaces below are declared here rather than in a shared port package.
// Go connects them structurally, so a caller supplies an adapter and this
// package never names one; declaring them at the point of use is also what
// keeps the port set honest -- it can only describe what this use case actually
// does, not what some other operation might want one day.
package resetcommands

import (
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// State is the command set and the way to persist a new one.
type State interface {
	// Commands returns the current command set.
	Commands() catalog.Catalog

	// Commit replaces the stored command set and records the operation with
	// the note that explains it.
	//
	// The note is rendered by the Notes port: the format is a storage detail,
	// and a use case that built it would own a second copy of it.
	Commit(commands catalog.Catalog, op operation.Operation, commandName, note string) error
}

// Output is the diagnostics stream. Nothing here is a result: `reset` prints no
// data a script would pipe.
type Output interface {
	Info(message string, data ...any)
	Table(commands catalog.Catalog) string
}

// Prompter asks the one question this use case has.
type Prompter interface {
	Confirm(question string) (bool, error)
}

// Notes renders the history note for the operation.
type Notes interface {
	Reset() string
}

// UseCase erases every managed command, after showing what will go.
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

// Run shows the command table, asks for confirmation and, on a yes, empties the
// command set.
//
// A declined confirmation is not a failure: nothing happened, and that was the
// answer, so the run reports it and returns nil.
func (u *UseCase) Run() error {
	u.output.Info("-----------------------------------------")
	u.output.Info(u.output.Table(u.state.Commands()))
	u.output.Info("")
	u.output.Info("All commands will be erased...")

	proceed, err := u.prompt.Confirm("Continue to reset")
	if err != nil {
		return err
	}

	if !proceed {
		u.output.Info("Cancelling reset...")
		return nil
	}

	u.output.Info("reset all commands")

	return u.state.Commit(catalog.New(nil, nil), operation.Reset, "", u.notes.Reset())
}
