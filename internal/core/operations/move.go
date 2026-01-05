package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/pkg/utils"
)

func Move(args []string, app *app.Data) {
	if len(args) == 0 {
		crerrors.PrintNotArgs("1 or more filepath/folderpath to move and output folder as last argument")
		return
	}

	lastIndex := len(args) - 1
	output := args[lastIndex]
	args = args[:lastIndex]
	for _, path := range args {
		var err error
		if utils.FilePathExists(path) {
			err = handlers.MoveFile(path, output, app)
		} else {
			err = handlers.MoveFolder(path, output, app)
		}

		if err != nil {
			fmt.Println(err)
		}
	}
}
