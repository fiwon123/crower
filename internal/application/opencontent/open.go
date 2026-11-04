// Package opencontent opens paths with the platform's default application, or
// the platform's environment editor.
//
// The operation is mostly orchestration: expand the arguments, hand each path
// to the Opener, report what ran, and record it. The defaults -- the
// configuration file for --file, the configuration folder for --folder -- come
// from the State port because the application does not know where crower keeps
// its files.
package opencontent

import (
	"fmt"
	"path/filepath"

	"github.com/fiwon123/crower/internal/domain/operation"
)

// Mode is what a request opens.
type Mode int

const (
	// ModePaths opens the given paths, glob-expanded.
	ModePaths Mode = iota

	// ModeFile opens files, defaulting to the configuration file.
	ModeFile

	// ModeFolder opens folders, defaulting to the configuration folder.
	ModeFolder

	// ModeSystem opens the platform's environment editor.
	ModeSystem
)

// Request is one open.
type Request struct {
	Mode Mode
	Args []string
}

// State is where the configuration lives and the history it is recorded in.
type State interface {
	// ConfigFile is the path of the configuration file.
	ConfigFile() string

	// ConfigFolder is the folder the configuration lives in.
	ConfigFolder() string

	// Record appends a history entry and persists the index.
	Record(op operation.Operation, commandName, note string) error
}

// Output is the diagnostics and the result stream.
type Output interface {
	Info(message string, data ...any)
	Result(text string)
}

// Opener opens one path with the platform's default application.
type Opener interface {
	// Open opens path and returns the command that ran, its trimmed output,
	// and any error.
	Open(path string) (command string, output string, err error)
}

// SystemEditor opens the platform's environment-variable editor.
//
// The method is not called Open because one adapter implements this port and
// the Opener, and a type cannot hold two methods named Open.
type SystemEditor interface {
	OpenEditor() error
}

// Notes renders the history note for each mode.
type Notes interface {
	OpenNote(args []string) string
	OpenFileNote(args []string) string
	OpenFolderNote(args []string) string
	OpenSystemNote() string
}

// UseCase runs one open.
type UseCase struct {
	state  State
	output Output
	opener Opener
	editor SystemEditor
	notes  Notes
}

// New builds the use case from its five dependencies.
func New(state State, output Output, opener Opener, editor SystemEditor, notes Notes) *UseCase {
	return &UseCase{state: state, output: output, opener: opener, editor: editor, notes: notes}
}

// Run opens what the request names and records it.
func (u *UseCase) Run(request Request) error {
	switch request.Mode {
	case ModePaths:
		return u.openPaths(request.Args)

	case ModeFile:
		return u.openFiles(request.Args)

	case ModeFolder:
		return u.openFolders(request.Args)

	case ModeSystem:
		return u.openSystem()

	default:
		return fmt.Errorf("unknown open mode %d", request.Mode)
	}
}

// openPaths opens every path, expanding each argument as a glob first.
//
// An argument that matches nothing is kept literally, which is how an open of a
// path the shell did not expand still works; a pattern that does not compile is
// skipped rather than failing the run.
func (u *UseCase) openPaths(args []string) error {
	paths := make([]string, 0, len(args))

	for _, arg := range args {
		matches, err := filepath.Glob(arg)
		if err != nil {
			continue
		}

		if len(matches) > 0 {
			paths = append(paths, matches...)
		} else {
			paths = append(paths, arg)
		}
	}

	if err := u.openAll(paths); err != nil {
		return err
	}

	return u.state.Record(operation.Open, "", u.notes.OpenNote(args))
}

// openFiles opens the given files, or the configuration file.
func (u *UseCase) openFiles(args []string) error {
	if len(args) == 0 {
		args = []string{u.state.ConfigFile()}
	}

	if err := u.openAll(args); err != nil {
		return err
	}

	return u.state.Record(operation.Open, "", u.notes.OpenFileNote(args))
}

// openFolders opens the given folders, or the configuration folder.
func (u *UseCase) openFolders(args []string) error {
	if len(args) == 0 {
		args = []string{u.state.ConfigFolder()}
	}

	if err := u.openAll(args); err != nil {
		return err
	}

	return u.state.Record(operation.Open, "", u.notes.OpenFolderNote(args))
}

// openSystem opens the platform's environment editor.
func (u *UseCase) openSystem() error {
	if err := u.editor.OpenEditor(); err != nil {
		return err
	}

	return u.state.Record(operation.Open, "", u.notes.OpenSystemNote())
}

// openAll opens every path, stopping at the first failure.
//
// Each path is attempted before the next: the loop used to return after the
// first one, so `crower open a b` silently opened only a.
func (u *UseCase) openAll(paths []string) error {
	for _, path := range paths {
		command, output, err := u.opener.Open(path)

		if command != "" {
			u.output.Info("opening: ", command)
		}

		if output != "" {
			u.output.Result(output)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
