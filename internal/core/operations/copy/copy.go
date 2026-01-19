package copyoperations

import (
	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func Copy(args []string, app *app.Data) {
	if len(args) == 0 {
		crerrors.PrintNotArgs("1 or more filepath/folderpath to copy and output folder as last argument", app)
		return
	}

	isCopyFile := false
	isCopyFolder := false
	lastIndex := len(args) - 1
	output := args[lastIndex]
	args = args[:lastIndex]
	for _, path := range args {
		var err error
		if utils.FilePathExists(path) {
			err = handlers.CopyFile(path, output, app)
			isCopyFile = true
		} else {
			err = handlers.CopyFolder(path, output, app)
			isCopyFolder = true
		}

		if err != nil {
			app.Logger.Info(err.Error())
		}
	}

	if isCopyFile && isCopyFolder {
		app.History.Add(state.Copy, notes.GenerateCopyNote(state.FileAndFolder, args))
	} else if isCopyFile {
		app.History.Add(state.Copy, notes.GenerateCopyNote(state.File, args))
	} else if isCopyFolder {
		app.History.Add(state.Copy, notes.GenerateCopyNote(state.Folder, args))
	}

	history.Save(app)
}
