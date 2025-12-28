package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func MoveFile(args []string, app *app.Data) {
	currentFolder := "./"
	fileName := ""
	destFolder := ""
	if len(args) > 2 {
		currentFolder = args[0]
		fileName = args[1]
		destFolder = args[2]
	} else {
		fmt.Println("needs to specify current foder path, file name and dest folder")
		return
	}

	handlers.MoveFile(currentFolder, fileName, destFolder, app)
}

func MoveFolder(args []string, app *app.Data) {
	currentFolder := "./"
	folderName := ""
	destFolder := ""
	if len(args) > 2 {
		currentFolder = args[0]
		folderName = args[1]
		destFolder = args[2]
	} else {
		fmt.Println("needs to specify current foder path, file name and dest folder")
		return
	}

	handlers.MoveFolder(currentFolder, folderName, destFolder, app)
}
