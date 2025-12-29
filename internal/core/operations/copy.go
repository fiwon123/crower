package operations

import (
	"fmt"

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
		fmt.Println("needs to specify file path and out folder")
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
		fmt.Println("needs to specify file path and out folder")
		return
	}

	handlers.CopyFolder(filePath, outFolder, app)
}
