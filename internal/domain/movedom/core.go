package movedom

import (
	stderrors "errors"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/errors"

	"github.com/fiwon123/crower/pkg/utils"
)

type Core struct {
	app     *app.Config
	handler *Handler
}

func NewCore(app *app.Config) *Core {

	handler := NewHandler(app)

	return &Core{
		handler: handler,
		app:     app,
	}
}

// Move moves every path to the last argument, the output folder.
//
// Like copy, every path is attempted and the failures are returned together.
func (c *Core) Move(args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("1 or more filepath/folderpath to move and output folder as last argument")
	}

	isMoveFile := false
	isMoveFolder := false
	lastIndex := len(args) - 1
	output := args[lastIndex]
	paths := args[:lastIndex]

	var failures []error

	// The note records the whole argument list, destination included: a
	// history entry that omits where the files went cannot describe what
	// happened.
	for _, path := range paths {
		var err error

		if utils.FilePathExists(path) {
			err = c.handler.MoveFile(path, output)
			isMoveFile = true
		} else {
			err = c.handler.MoveFolder(path, output)
			isMoveFolder = true
		}

		if err != nil {
			failures = append(failures, err)
		}
	}

	if len(failures) > 0 {
		return stderrors.Join(failures...)
	}

	if isMoveFile && isMoveFolder {
		c.app.History.Add(operations.Move, newMoveNote(operations.FileAndFolder, args))
	} else if isMoveFile {
		c.app.History.Add(operations.Move, newMoveNote(operations.File, args))
	} else if isMoveFolder {
		c.app.History.Add(operations.Move, newMoveNote(operations.Folder, args))
	}

	return c.app.Save()
}
