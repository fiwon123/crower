package moveoperations

import (
	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func Move(args []string, app *app.Data) {
	if len(args) == 0 {
		crerrors.PrintNotArgs("1 or more filepath/folderpath to move and output folder as last argument", app)
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
			err = handlers.MoveFile(path, output, app)
			isMoveFile = true
		} else {
			err = handlers.MoveFolder(path, output, app)
			isMoveFolder = true
		}

		if err != nil {
			app.Logger.Error(err.Error())
			return
		}
	}

	if isMoveFile && isMoveFolder {
		app.History.Add(state.Move, notes.GenerateMoveNote(state.FileAndFolder, args))
	} else if isMoveFile {
		app.History.Add(state.Move, notes.GenerateMoveNote(state.File, args))
	} else if isMoveFolder {
		app.History.Add(state.Move, notes.GenerateMoveNote(state.Folder, args))
	}

	history.Save(app)
}
