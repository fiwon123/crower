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

	// Commit replaces the stored command set and records the operation that
	// explains it.
	//
	// The use case describes what happened -- the operation and the crower
	// subcommand that reproduces it -- rather than how it is written: the
	// storage format is the adapter's business.
	Commit(commands catalog.Catalog, op operation.Operation, replay string) error
}

// Output is the diagnostics stream. Nothing here is a result: `reset` prints no
// data a script would pipe.
type Output interface {
	Info(message string)
	Table(commands catalog.Catalog) string
}

// Prompter asks the one question this use case has.
type Prompter interface {
	Confirm(question string) (bool, error)
}

// UseCase erases every managed command, after showing what will go.
type UseCase struct {
	state  State
	output Output
	prompt Prompter
}

// New builds the use case from its three dependencies.
func New(state State, output Output, prompt Prompter) *UseCase {
	return &UseCase{state: state, output: output, prompt: prompt}
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

	return u.state.Commit(catalog.New(nil, nil), operation.Reset, "reset")
}
