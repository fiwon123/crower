package copydom

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

// Copy copies every path to the last argument, the output folder.
//
// Every path is attempted even when one fails, and the failures are returned
// together: a batch that copied seven of nine files must say so rather than
// report success for the whole run.
func (c *Core) Copy(args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("1 or more filepath/folderpath to copy and output folder as last argument")
	}

	isCopyFile := false
	isCopyFolder := false
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
			err = c.handler.CopyFile(path, output)
			isCopyFile = true
		} else {
			err = c.handler.CopyFolder(path, output)
			isCopyFolder = true
		}

		if err != nil {
			failures = append(failures, err)
		}
	}

	if len(failures) > 0 {
		return stderrors.Join(failures...)
	}

	if isCopyFile && isCopyFolder {
		c.app.History.Add(operations.Copy, newCopyNote(operations.FileAndFolder, args))
	} else if isCopyFile {
		c.app.History.Add(operations.Copy, newCopyNote(operations.File, args))
	} else if isCopyFolder {
		c.app.History.Add(operations.Copy, newCopyNote(operations.Folder, args))
	}

	return c.app.Save()
}
