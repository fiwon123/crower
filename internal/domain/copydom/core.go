package copydom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/errors"

	"github.com/fiwon123/crower/pkg/crowerutils"
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

func (c *Core) Copy(args []string) {
	if len(args) == 0 {
		errors.PrintNotArgs("1 or more filepath/folderpath to copy and output folder as last argument", c.app)
		return
	}

	isCopyFile := false
	isCopyFolder := false
	lastIndex := len(args) - 1
	output := args[lastIndex]
	args = args[:lastIndex]
	for _, path := range args {
		var err error
		if crowerutils.FilePathExists(path) {
			err = c.handler.CopyFile(path, output)
			isCopyFile = true
		} else {
			err = c.handler.CopyFolder(path, output)
			isCopyFolder = true
		}

		if err != nil {
			c.app.Logger.Info(err.Error())
		}
	}

	if isCopyFile && isCopyFolder {
		c.app.History.Add(operations.Copy, newCopyNote(operations.FileAndFolder, args))
	} else if isCopyFile {
		c.app.History.Add(operations.Copy, newCopyNote(operations.File, args))
	} else if isCopyFolder {
		c.app.History.Add(operations.Copy, newCopyNote(operations.Folder, args))
	}

	c.app.Save()
}
