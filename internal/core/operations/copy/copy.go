package copyoperations

import (
	"github.com/fiwon123/crower/internal/crowererrors"
	appdata "github.com/fiwon123/crower/internal/data/app"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	copyhandlers "github.com/fiwon123/crower/internal/handlers/copy"

	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

func Copy(args []string, app *appdata.Data) {
	if len(args) == 0 {
		crowererrors.PrintNotArgs("1 or more filepath/folderpath to copy and output folder as last argument", app)
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
			err = copyhandlers.CopyFile(path, output, app)
			isCopyFile = true
		} else {
			err = copyhandlers.CopyFolder(path, output, app)
			isCopyFolder = true
		}

		if err != nil {
			app.Logger.Info(err.Error())
		}
	}

	if isCopyFile && isCopyFolder {
		app.History.Add(operationsdata.Copy, notes.GenerateCopyNote(operationsdata.FileAndFolder, args))
	} else if isCopyFile {
		app.History.Add(operationsdata.Copy, notes.GenerateCopyNote(operationsdata.File, args))
	} else if isCopyFolder {
		app.History.Add(operationsdata.Copy, notes.GenerateCopyNote(operationsdata.Folder, args))
	}

	history.Save(app)
}
