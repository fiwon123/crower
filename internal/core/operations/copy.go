package operations

import (
	"github.com/fiwon123/crower/internal/cterrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func CopyFile(args []string, app *app.Data) {
	filePath := ""
	outFolder := ""
	if len(args) > 1 {
		filePath = args[0]
		outFolder = args[1]
	} else {
		cterrors.PrintNotFileAndOutputPath()
		return
	}

	handlers.CopyFile(filePath, outFolder, app)
}

func CopyFolder(args []string, app *app.Data) {
	filePath := ""
	outFolder := ""
	if len(args) > 1 {
		filePath = args[0]
		outFolder = args[1]
	} else {
		cterrors.PrintNotFileAndOutputPath()
		return
	}

	handlers.CopyFolder(filePath, outFolder, app)
}
