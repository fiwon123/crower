package moveoperations

import (
	"github.com/fiwon123/crower/internal/crowererrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	movehandlers "github.com/fiwon123/crower/internal/handlers/move"

	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

func Move(args []string, app *app.Data) {
	if len(args) == 0 {
		crowererrors.PrintNotArgs("1 or more filepath/folderpath to move and output folder as last argument", app)
		return
	}

	isMoveFile := false
	isMoveFolder := false
	lastIndex := len(args) - 1
	output := args[lastIndex]
	args = args[:lastIndex]
	for _, path := range args {
		var err error
		if crowerutils.FilePathExists(path) {
			err = movehandlers.MoveFile(path, output, app)
			isMoveFile = true
		} else {
			err = movehandlers.MoveFolder(path, output, app)
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
