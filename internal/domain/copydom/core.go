package copydom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/crowererrors"
	copynotesdata "github.com/fiwon123/crower/internal/data/notes/copy"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"

	historyhelper "github.com/fiwon123/crower/internal/helper/history"
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
		crowererrors.PrintNotArgs("1 or more filepath/folderpath to copy and output folder as last argument", c.app)
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
		c.app.History.Add(operationsdata.Copy, copynotesdata.NewCopyNote(operationsdata.FileAndFolder, args))
	} else if isCopyFile {
		c.app.History.Add(operationsdata.Copy, copynotesdata.NewCopyNote(operationsdata.File, args))
	} else if isCopyFolder {
		c.app.History.Add(operationsdata.Copy, copynotesdata.NewCopyNote(operationsdata.Folder, args))
	}

	historyhelper.Save(c.app)
}
