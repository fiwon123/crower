package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func Extract(args []string, app *app.Data) {
	folderPath := "./"
	fileName := ""
	outDir := folderPath
	if len(args) > 1 {
		folderPath = args[0]
		fileName = args[1]

		if len(args) > 2 {
			outDir = args[2]
		} else {
			outDir = folderPath
		}
	} else {
		fmt.Println("at least pass 3 parameters, folderPath, fileName, outDir")
		return
	}

	handlers.Extract(folderPath, fileName, outDir, app)
}
