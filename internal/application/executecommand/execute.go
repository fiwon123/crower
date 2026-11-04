// Package executecommand runs a managed command.
//
// The rules are the operation's own: how parameters fill a line, how a Steps
// command runs its steps, and when a step is a nested crower call that needs the
// depth guard. The technology -- starting a shell, starting a child process --
// is behind the Runner port.
//
// The same package serves the command line and the two user interfaces. The CLI
// path (Run) prompts, reports and records; the UI path (Execute) answers with the
// output and the command name, because those callers own their own interaction
// and have never written a history entry.
package executecommand

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/errors"
)

// Request is one execution.
type Request struct {
	// Key is the command to run, by name or alias. Empty means "ask".
	Key string

	// Params are the positional values applied to the command.
	Params []string

	// Last, when set, runs the command the last operation of this kind
	// touched. The --last forms take no key and ask nothing.
	Last operation.Operation
}

// State is the command set and the history it is recorded in.
type State interface {
	// Commands returns the current command set.
	Commands() catalog.Catalog

	// LastCommand returns the command name the last operation of this kind
	// touched.
	LastCommand(op operation.Operation) (string, bool)

	// Record appends a history entry and persists the index.
	Record(op operation.Operation, commandName, note string) error
}

// Output is the diagnostics stream, the result stream, and the debug log.
type Output interface {
	// Info writes a diagnostic line.
	Info(message string, data ...any)

	// Table renders the current command table.
	Table(commands catalog.Catalog) string

	// Result writes a command's own output: data a caller may pipe, not a
	// diagnostic about the run.
	Result(text string)

	// Debug writes a line only the debug log keeps.
	Debug(message string, data ...any)
}

// Prompter asks the questions this use case has.
type Prompter interface {
	ChooseCommand() (string, error)
	Confirm(question string) (bool, error)
}

// Runner runs shell lines.
type Runner interface {
	// Run runs one line and returns its combined output.
	Run(line string) (string, error)

	// RunNested runs a line that invokes crower again, refusing to go past the
	// nesting limit.
	RunNested(line string) (string, error)
}

// Notes renders the history note for the execution.
type Notes interface {
	ExecuteCommandNote(command catalog.Command, params []string) string
	ExecuteLastNote(op operation.Operation, command catalog.Command, args []string) string
}

// UseCase runs one command.
type UseCase struct {
	state  State
	output Output
	prompt Prompter
	runner Runner
	notes  Notes
}

// New builds the use case from its five dependencies.
func New(state State, output Output, prompt Prompter, runner Runner, notes Notes) *UseCase {
	return &UseCase{state: state, output: output, prompt: prompt, runner: runner, notes: notes}
}

// Run runs what the request names, prompting for the command when the request
// names none, and records the operation.
func (u *UseCase) Run(request Request) error {
	if request.Last != "" {
		return u.runLast(request)
	}

	return u.runCommand(request)
}

// Execute runs the command a key resolves to, without prompting and without
// recording history.
//
// It returns the output and the command's name for the caller to show; the
// terminal and desktop interfaces use it because they own their own interaction
// and have never added history entries.
func (u *UseCase) Execute(key string, params []string) (string, string, error) {
	output, command, err := u.execute(key, params)

	return output, command.Name, err
}

// runCommand runs a positional command, asking which one when the request names
// none.
func (u *UseCase) runCommand(request Request) error {
	key := request.Key

	if key == "" {
		chosen, err := u.chooseCommand()
		if err != nil {
			return err
		}

		key = chosen

		if _, found := u.state.Commands().Resolve(key); !found {
			u.output.Info("Command not found, try to select one.")

			chosen, err := u.chooseCommand()
			if err != nil {
				return err
			}

			key = chosen

			if _, found = u.state.Commands().Resolve(key); !found {
				return errors.CommandNotFound(key)
			}
		}

		command, _ := u.state.Commands().Resolve(key)
		u.showCommand(command)

		proceed, err := u.prompt.Confirm("Continue to execute")
		if err != nil {
			return err
		}

		if !proceed {
			u.output.Info("Cancelling execute...")
			return nil
		}
	}

	output, executed, err := u.execute(key, request.Params)
	if err != nil {
		return err
	}

	u.report(output)

	return u.state.Record(operation.Execute, executed.Name, u.notes.ExecuteCommandNote(executed, request.Params))
}

// runLast runs the command the last operation of a kind touched.
//
// The --last forms take no key and ask nothing: the flags describe the whole
// request, and a prompt in the middle of a scripted call would hang it.
func (u *UseCase) runLast(request Request) error {
	name, found := u.state.LastCommand(request.Last)
	if !found {
		return errors.CommandNotFound(string(request.Last))
	}

	output, executed, err := u.execute(name, request.Params)
	if err != nil {
		return err
	}

	u.report(output)

	return u.state.Record(operation.Execute, executed.Name, u.notes.ExecuteLastNote(request.Last, executed, request.Params))
}

// execute resolves the key and runs the command's steps.
func (u *UseCase) execute(key string, params []string) (string, catalog.Command, error) {
	command, found := u.state.Commands().Resolve(key)
	if !found {
		return "", catalog.Command{}, errors.CommandNotFound(key)
	}

	output, err := u.runSteps(command, params)

	return output, command, err
}

// runSteps applies the parameters and runs every step in order.
//
// A command with steps never runs its exec, so an unusable brace there must not
// fail a run that never touches it.
func (u *UseCase) runSteps(command catalog.Command, params []string) (string, error) {
	steps := command.Steps

	if len(steps) == 0 {
		execLine, err := applyParams(command.Exec, params)
		if err != nil {
			return "", err
		}

		u.output.Debug("exec: " + execLine)

		steps = []string{execLine}
	} else {
		substituted, err := applyParamsToSteps(steps, params)
		if err != nil {
			return "", err
		}

		steps = substituted
	}

	var allOut strings.Builder

	for i, step := range steps {
		step = strings.TrimSpace(step)
		if step == "" {
			continue
		}

		if isNestedCrower(step) {
			out, err := u.runner.RunNested(step)
			allOut.WriteString(out)

			if err != nil {
				return allOut.String(), err
			}

			continue
		}

		u.output.Debug(fmt.Sprintf("[%d/%d] %s", i+1, len(steps), step))

		out, err := u.runner.Run(step)
		allOut.WriteString(out)

		if err != nil {
			return allOut.String(), err
		}
	}

	return allOut.String(), nil
}

// chooseCommand prints the table and asks for a row.
func (u *UseCase) chooseCommand() (string, error) {
	u.output.Info(u.output.Table(u.state.Commands()))

	key, err := u.prompt.ChooseCommand()
	if err != nil {
		return "", err
	}

	if key == "" {
		return "", errors.InvalidRows()
	}

	return key, nil
}

// showCommand prints the fields the confirmation is about.
func (u *UseCase) showCommand(command catalog.Command) {
	u.output.Info("-----------------------------------------")
	u.output.Info("Name:    ", command.Name)
	u.output.Info("Aliases: ", command.Aliases)
	u.output.Info("Exec:    ", command.Exec)
	u.output.Info("")
}

// report writes a command's own output to the result stream.
//
// The output belongs to the executed command, so it is data a caller may pipe,
// not a diagnostic about the run. Printing it through the diagnostic writer put
// it on stderr and appended a space.
func (u *UseCase) report(output string) {
	trimmed := strings.TrimRight(output, "\n")
	if trimmed == "" {
		return
	}

	u.output.Result(trimmed)
}
