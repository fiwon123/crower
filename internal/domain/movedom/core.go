package movedom

import (
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

func (c *Core) Move(args []string) {
	if len(args) == 0 {
		errors.PrintNotArgs("1 or more filepath/folderpath to move and output folder as last argument", c.app)
		return
	}

	isMoveFile := false
	isMoveFolder := false
	lastIndex := len(args) - 1
	output := args[lastIndex]
	args = args[:lastIndex]
	for _, path := range args {
		var err error
		if utils.FilePathExists(path) {
			err = c.handler.MoveFile(path, output)
			isMoveFile = true
		} else {
			err = c.handler.MoveFolder(path, output)
			isMoveFolder = true
		}

		if err != nil {
			c.app.Logger.Error(err.Error())
			return
		}
	}

	if isMoveFile && isMoveFolder {
		c.app.History.Add(operations.Move, newMoveNote(operations.FileAndFolder, args))
	} else if isMoveFile {
		c.app.History.Add(operations.Move, newMoveNote(operations.File, args))
	} else if isMoveFolder {
		c.app.History.Add(operations.Move, newMoveNote(operations.Folder, args))
	}

	c.app.Save()
}
