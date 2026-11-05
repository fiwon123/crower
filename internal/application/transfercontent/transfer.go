// Package transfercontent copies or moves paths to a destination folder.
//
// `copy` and `move` are the same flow with one word changed, so they share this
// package and a Request that names the operation. The platform command is
// behind the Transfer port; the note's sub-operation follows what was actually
// transferred.
package transfercontent

import (
	"context"
	stderrors "errors"

	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/errors"
)

// Request is one copy or move.
type Request struct {
	// Op is operation.Copy or operation.Move.
	Op operation.Operation

	// Args are the positional values: one or more sources, then the
	// destination folder as the last one.
	Args []string

	// Context stops the batch between sources. The platform command already
	// running is not interrupted; no further source starts. Nil means "run to
	// the end".
	Context context.Context

	// Progress is told how many sources are done, out of how many, and which
	// one is starting. It may be nil.
	Progress func(done, total int, current string)
}

// State is the history the operation is recorded in.
type State interface {
	// Record appends a history entry and persists the index.
	Record(op operation.Operation, commandName, note string) error
}

// Files answers whether a source is a file or a folder.
type Files interface {
	// IsFile reports whether path exists and is a file. Anything else --
	// a folder, or a path that does not exist -- takes the folder branch,
	// which is what this command has always done.
	IsFile(path string) bool
}

// Transfer copies and moves paths.
type Transfer interface {
	Copy(source, destination string, folder bool) error
	Move(source, destination string) error
}

// Notes renders the history note for the operation.
type Notes interface {
	CopyNote(subOp operation.SubOperation, args []string) string
	MoveNote(subOp operation.SubOperation, args []string) string
}

// UseCase runs one copy or move.
type UseCase struct {
	state    State
	files    Files
	transfer Transfer
	notes    Notes
}

// New builds the use case from its four dependencies.
func New(state State, files Files, transfer Transfer, notes Notes) *UseCase {
	return &UseCase{state: state, files: files, transfer: transfer, notes: notes}
}

// Run copies or moves every source to the destination.
//
// Every source is attempted even when one fails, and the failures are returned
// together: a batch that copied seven of nine files must say so rather than
// report success for the whole run.
func (u *UseCase) Run(request Request) error {
	verb := "copy"
	if request.Op == operation.Move {
		verb = "move"
	}

	// One argument names a source and no destination. It used to "succeed" by
	// doing nothing; there is nothing to do it to, so it is refused.
	if len(request.Args) < 2 {
		return errors.NotEnoughArgs("1 or more filepath/folderpath to " + verb + " and output folder as last argument")
	}

	lastIndex := len(request.Args) - 1
	destination := request.Args[lastIndex]
	paths := request.Args[:lastIndex]

	ctx := request.Context
	if ctx == nil {
		ctx = context.Background()
	}

	sawFile := false
	sawFolder := false

	var failures []error

	for index, path := range paths {
		if err := ctx.Err(); err != nil {
			return errors.Cancelled(index, len(paths))
		}

		if request.Progress != nil {
			request.Progress(index, len(paths), path)
		}

		// The callback is where a window's Cancel button is pressed, so the
		// check runs again: a cancel asked for on this item skips it.
		if err := ctx.Err(); err != nil {
			return errors.Cancelled(index, len(paths))
		}

		folder := !u.files.IsFile(path)

		if folder {
			sawFolder = true
		} else {
			sawFile = true
		}

		var err error

		if request.Op == operation.Move {
			err = u.transfer.Move(path, destination)
		} else {
			err = u.transfer.Copy(path, destination, folder)
		}

		if err != nil {
			failures = append(failures, err)
		}
	}

	if len(failures) > 0 {
		return stderrors.Join(failures...)
	}

	subOp := subOperation(sawFile, sawFolder)

	if request.Op == operation.Move {
		return u.state.Record(operation.Move, "", u.notes.MoveNote(subOp, request.Args))
	}

	return u.state.Record(operation.Copy, "", u.notes.CopyNote(subOp, request.Args))
}

// subOperation names what was transferred.
func subOperation(file, folder bool) operation.SubOperation {
	switch {
	case file && folder:
		return operation.FileAndFolder
	case folder:
		return operation.Folder
	default:
		return operation.File
	}
}
