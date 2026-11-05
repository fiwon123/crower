// Package createcommand adds a managed command, a file, a folder, an
// environment variable or a PATH entry.
//
// `create` looks like one operation with eight flags, but underneath it is five
// jobs: building a command (from the line, from a running process, from a
// script file, or from the last command that ran), creating paths, editing the
// environment, and opening the folder a new script lives in. They share this
// package because they share the CLI and the history operation, and each mode
// has its own narrow ports.
//
// The command layer decides which mode a run is in -- it owns the flags -- and
// this package owns what that mode asks, writes and records.
package createcommand

import (
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/errors"
)

// Mode is what a request creates.
type Mode int

const (
	// ModeCommand adds a managed command from a name and the line it runs.
	ModeCommand Mode = iota

	// ModeProcess adds one built from a running process.
	ModeProcess

	// ModeScript writes a script file and registers it.
	ModeScript

	// ModeLast rebuilds the command the last operation of Last ran.
	ModeLast

	// ModeFile creates one or more files.
	ModeFile

	// ModeFolder creates one or more folders.
	ModeFolder

	// ModeSystemVariable adds one environment variable.
	ModeSystemVariable

	// ModeSystemPath adds one PATH entry.
	ModeSystemPath
)

// Request is one create.
type Request struct {
	// Mode is which job to run.
	Mode Mode

	// Name is the new command's name for ModeProcess, where the flag carries
	// it. Empty means the first positional value is the name, as before.
	Name string

	// Aliases are the aliases supplied on the command line.
	Aliases []string

	// Icon and Group are the display metadata supplied on the command line.
	Icon  string
	Group string

	// Description, Params, Cwd, Shell and Timeout are the optional execution
	// metadata supplied on the command line.
	Description string
	Params      []string
	Cwd         string
	Shell       string
	Timeout     string

	// Args are the positional values the command line carried.
	Args []string

	// Last names the operation whose command is rebuilt by ModeLast.
	Last operation.Operation
}

// LastRun is what the history recorded about the last operation of a kind.
type LastRun struct {
	// CommandName is the command the operation touched, when it touched one.
	CommandName string

	// Exec is the line the operation ran, as the note recorded it.
	Exec string

	// Steps are the lines a Steps command ran, and HasSteps says whether the
	// note (or the catalog, for older notes) recorded any.
	Steps    []string
	HasSteps bool
}

// State is the command set, the history it is recorded in, and the way to
// persist both.
type State interface {
	// Commands returns the current command set.
	Commands() catalog.Catalog

	// Save replaces the stored command set.
	//
	// It is separate from Record -- rather than the single Commit other use
	// cases use -- because `create --script` opens the scripts folder before
	// the history entry is written, which is the order it has always had.
	Save(commands catalog.Catalog) error

	// Record appends a history entry and persists it.
	Record(op operation.Operation, commandName, note string) error

	// LastOperation describes the last operation of a kind: the command it
	// touched, the line it ran, and whether it ran steps.
	LastOperation(op operation.Operation) (LastRun, bool)
}

// Output is the diagnostics stream.
type Output interface {
	Info(message string, data ...any)
}

// Prompter asks for the fields a create was not given.
//
// The question texts belong to the prompter: they are the ones crower has
// always shown, and they differ from the update prompts.
type Prompter interface {
	CreateName() (string, error)
	CreateExec() (string, error)
	NewAliases() ([]string, error)
	Confirm(question string) (bool, error)
}

// Creator creates files and folders with the host's own tools.
type Creator interface {
	File(path string) error
	Folder(path string) error
}

// Scripts writes the script file a `create --script` command runs.
type Scripts interface {
	// Create makes the script file for name and returns its path and the
	// folder that holds it.
	Create(name string) (path string, folder string, err error)
}

// Processes resolves a process name or pid to the line that starts it again.
type Processes interface {
	ExecLine(process string) (string, error)
}

// Environment edits the variables the user's shell reads at start.
type Environment interface {
	SetVariable(name, value string) (string, error)
	AddPath(path string) (string, error)
}

// Opener opens a path with the host's file manager, so the user can edit the
// script that was just created.
type Opener interface {
	Open(path string) error
}

// Notes renders the history note for each mode.
//
// The names carry the Note suffix because one adapter implements every port,
// and the Scripts and Creator ports already own CreateScript, CreateFile and
// CreateFolder.
type Notes interface {
	CreateCommandNote(command catalog.Command) string
	CreateScriptNote(command catalog.Command) string
	CreateLastNote(command catalog.Command) string
	CreateProcessNote(command catalog.Command, args []string) string
	CreateSystemVariableNote(args []string) string
	CreateSystemPathNote(args []string) string
	CreateFileNote(args []string) string
	CreateFolderNote(args []string) string
}

// UseCase runs one create.
type UseCase struct {
	state       State
	output      Output
	prompt      Prompter
	creator     Creator
	scripts     Scripts
	processes   Processes
	environment Environment
	opener      Opener
	notes       Notes
}

// New builds the use case from its nine dependencies.
func New(
	state State,
	output Output,
	prompt Prompter,
	creator Creator,
	scripts Scripts,
	processes Processes,
	environment Environment,
	opener Opener,
	notes Notes,
) *UseCase {
	return &UseCase{
		state:       state,
		output:      output,
		prompt:      prompt,
		creator:     creator,
		scripts:     scripts,
		processes:   processes,
		environment: environment,
		opener:      opener,
		notes:       notes,
	}
}

// Run creates what the request names.
func (u *UseCase) Run(request Request) error {
	switch request.Mode {
	case ModeProcess:
		return u.createProcess(request)

	case ModeScript:
		return u.createScript(request)

	case ModeLast:
		return u.createLast(request)

	case ModeFile:
		return u.createFile(request)

	case ModeFolder:
		return u.createFolder(request)

	case ModeSystemVariable:
		return u.createSystemVariable(request)

	case ModeSystemPath:
		return u.createSystemPath(request)

	case ModeCommand:
		return u.createCommand(request)

	default:
		// A mode this build does not know must not fall back to the branch
		// that adds a command.
		return fmt.Errorf("unknown create mode %d", request.Mode)
	}
}

// createCommand adds a managed command, asking for the fields the command line
// did not carry.
//
// Two positional values are the whole command and skip the questions; any
// other count falls back to asking for whatever is missing, which is the rule
// create has always had.
func (u *UseCase) createCommand(request Request) error {
	name, exec := "", ""

	if len(request.Args) == 2 {
		name, exec = request.Args[0], request.Args[1]
	} else {
		var err error

		if name == "" {
			if name, err = u.prompt.CreateName(); err != nil {
				return err
			}
		}

		if len(request.Aliases) == 0 {
			addAliases, err := u.prompt.Confirm("Do you want to add alias")
			if err != nil {
				return err
			}

			if addAliases {
				if request.Aliases, err = u.prompt.NewAliases(); err != nil {
					return err
				}
			}
		}

		if exec == "" {
			if exec, err = u.prompt.CreateExec(); err != nil {
				return err
			}
		}
	}

	commands, created, err := u.add(request.command(name, exec))
	if err != nil {
		return fmt.Errorf("add command %q: %w", name, err)
	}

	if err := u.state.Save(commands); err != nil {
		return err
	}

	u.output.Info("added new command: ", name)

	return u.state.Record(operation.Create, created.Name, u.notes.CreateCommandNote(created))
}

// createProcess adds a command built from a running process.
func (u *UseCase) createProcess(request Request) error {
	name := request.Name
	args := request.Args

	if len(args) > 0 && name == "" {
		name = args[0]
		args = args[1:]
	}

	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return errors.NotEnoughArgs("process name or pid")
	}

	execLine, err := u.processes.ExecLine(args[0])
	if err != nil {
		return fmt.Errorf("add command by process: %w", err)
	}

	commands, created, err := u.add(request.command(name, execLine))
	if err != nil {
		return fmt.Errorf("add command by process: %w", err)
	}

	if err := u.state.Save(commands); err != nil {
		return err
	}

	u.output.Info("added new command by process: ", name)

	return u.state.Record(operation.Create, created.Name, u.notes.CreateProcessNote(created, request.Args))
}

// createScript writes a script file, registers it, and opens its folder.
//
// The folder is opened before the entry is recorded, which is the order it has
// always had: the user is sent to the script to write it, and the history
// entry exists only once that step succeeded.
func (u *UseCase) createScript(request Request) error {
	if len(request.Args) == 0 {
		return errors.NotEnoughArgs("name")
	}

	name := request.Args[0]

	path, folder, err := u.scripts.Create(name)
	if err != nil {
		return err
	}

	commands, created, err := u.add(request.command(name, path))
	if err != nil {
		return fmt.Errorf("add command %q: %w", name, err)
	}

	if err := u.state.Save(commands); err != nil {
		return err
	}

	u.output.Info("added new command: ", name)

	if err := u.opener.Open(folder); err != nil {
		return err
	}

	return u.state.Record(operation.Create, created.Name, u.notes.CreateScriptNote(created))
}

// createLast rebuilds the command the last operation of a kind ran.
func (u *UseCase) createLast(request Request) error {
	if len(request.Args) == 0 {
		return errors.NotEnoughArgs("name")
	}

	if len(request.Args) > 1 {
		return errors.TooManyArgs("one name")
	}

	name := request.Args[0]

	run, found := u.state.LastOperation(request.Last)
	if !found {
		return errors.CommandNotFound(string(request.Last))
	}

	// A command with steps never runs its exec, so reusing that value would
	// rebuild a command that does something else. Say so instead.
	if run.HasSteps {
		if len(run.Steps) == 0 {
			return fmt.Errorf("the last %q operation (%q) has steps that cannot be read; create --execute cannot rebuild them", request.Last, run.CommandName)
		}

		return fmt.Errorf("the last %q operation (%q) ran %d step(s); create --execute cannot rebuild them", request.Last, run.CommandName, len(run.Steps))
	}

	if run.Exec == "" {
		return fmt.Errorf("the last %q operation carries no exec to reuse", request.Last)
	}

	commands, created, err := u.add(request.command(name, run.Exec))
	if err != nil {
		return fmt.Errorf("add command %q: %w", name, err)
	}

	if err := u.state.Save(commands); err != nil {
		return err
	}

	u.output.Info("added new command: ", name)

	return u.state.Record(operation.Create, created.Name, u.notes.CreateLastNote(created))
}

// createFile creates every path, reporting all failures together.
func (u *UseCase) createFile(request Request) error {
	if len(request.Args) == 0 {
		return errors.NotEnoughArgs("file path")
	}

	var failures []error

	for _, path := range request.Args {
		if err := u.creator.File(path); err != nil {
			failures = append(failures, err)
		}
	}

	if len(failures) > 0 {
		return stderrors.Join(failures...)
	}

	return u.state.Record(operation.Create, "", u.notes.CreateFileNote(request.Args))
}

// createFolder creates every path, reporting all failures together.
func (u *UseCase) createFolder(request Request) error {
	if len(request.Args) == 0 {
		return errors.NotEnoughArgs("folder path")
	}

	var failures []error

	for _, path := range request.Args {
		if err := u.creator.Folder(path); err != nil {
			failures = append(failures, err)
		}
	}

	if len(failures) > 0 {
		return stderrors.Join(failures...)
	}

	return u.state.Record(operation.Create, "", u.notes.CreateFolderNote(request.Args))
}

// createSystemVariable adds one environment variable and records it.
func (u *UseCase) createSystemVariable(request Request) error {
	if len(request.Args) < 2 {
		return errors.NotEnoughArgs("var name and var value")
	}

	out, err := u.environment.SetVariable(request.Args[0], request.Args[1])
	if err != nil {
		return err
	}

	u.output.Info(out)

	return u.state.Record(operation.Create, "", u.notes.CreateSystemVariableNote(request.Args))
}

// createSystemPath adds one PATH entry and records it.
func (u *UseCase) createSystemPath(request Request) error {
	if len(request.Args) == 0 {
		return errors.NotEnoughArgs("path")
	}

	out, err := u.environment.AddPath(request.Args[0])
	if err != nil {
		return err
	}

	u.output.Info(out)

	return u.state.Record(operation.Create, "", u.notes.CreateSystemPathNote(request.Args))
}

// add checks the rules and returns the catalog that holds the new command.
//
// It does not persist anything: the caller decides when the command file is
// written, because `--script` has a step in between.
func (u *UseCase) add(command catalog.Command) (catalog.Catalog, catalog.Command, error) {
	commands, err := u.state.Commands().Add(command)
	if err != nil {
		return catalog.Catalog{}, catalog.Command{}, err
	}

	created, _ := commands.Get(command.Name)

	return commands, created, nil
}

// command is the command a create mode adds: the line it runs plus the metadata
// the flags supplied. Every mode that ends in a command builds it here, so the
// metadata cannot reach one mode and miss another.
func (r Request) command(name, exec string) catalog.Command {
	return catalog.Command{
		Name:        name,
		Aliases:     r.Aliases,
		Exec:        exec,
		Icon:        r.Icon,
		Group:       r.Group,
		Description: r.Description,
		Params:      r.Params,
		Cwd:         r.Cwd,
		Shell:       r.Shell,
		Timeout:     r.Timeout,
	}
}
