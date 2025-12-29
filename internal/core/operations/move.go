package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func MoveFile(args []string, app *app.Data) {
	filePath := ""
	destFolder := ""
	if len(args) > 1 {
		filePath = args[0]
		destFolder = args[1]
	} else {
		fmt.Println("needs to specify file path and  output folder")
		return
	}

	handlers.MoveFile(filePath, destFolder, app)
}

func MoveFolder(args []string, app *app.Data) {
	folderPath := ""
	destFolder := ""
	if len(args) > 1 {
		folderPath = args[0]
		destFolder = args[1]
	} else {
		fmt.Println("needs to specify folder path and  output folder")
		return
	}

	handlers.MoveFolder(folderPath, destFolder, app)
}
