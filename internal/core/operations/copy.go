package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/cterrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/pkg/utils"
)

func Copy(args []string, app *app.Data) {
	if len(args) == 0 {
		cterrors.PrintNotArgs("1 or more filepath/folderpath to copy and output folder as last argument")
		return
	}

	lastIndex := len(args) - 1
	output := args[lastIndex]
	args = args[:lastIndex]
	for _, path := range args {
		var err error
		if utils.IsValidFilePath(path) {
			err = handlers.CopyFile(path, output, app)
		} else {
			err = handlers.CopyFolder(path, output, app)
		}

		if err != nil {
			fmt.Println(err)
		}
	}
}
