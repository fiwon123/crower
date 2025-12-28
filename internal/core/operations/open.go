package operations

import (
	"fmt"
	"path"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func Open(args []string, app *app.Data) {
	currentPath := app.CfgFilePath
	if len(args) == 0 {
		args = append(args, currentPath)
	}
	handlers.Open(args, app)
}

func OpenFolder(args []string, app *app.Data) {
	currentPath := app.CfgFilePath
	if len(args) == 0 {
		args = append(args, path.Dir(currentPath))
	}

	handlers.Open(args, app)
}

func OpenSystem(app *app.Data) {
	out, err := handlers.OpenSystem(app)
	if err != nil {
		fmt.Println("failed to open system variable: ", err)
		return
	}

	fmt.Print(string(out))
}
