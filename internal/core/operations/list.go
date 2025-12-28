package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/payload"
	"github.com/fiwon123/crower/internal/handlers"
)

func ListCommands(app *app.Data) {
	handlers.ListCommands(app)
}

func ListProcess(payload payload.Data, app *app.Data) {
	handlers.ListProcess(payload.Args, app)
}

func ListHistory(app *app.Data) {
	handlers.ListHistory(app)
}

func ListFolder(args []string, app *app.Data) {
	currentPath := "./"
	if len(args) > 0 {
		currentPath = args[0]
	}

	out, err := handlers.ListFolder(currentPath, app)
	if err != nil {
		fmt.Println("failed to list folder: ", err)
		return
	}

	fmt.Print(string(out))
}

func ListSystem(app *app.Data) {
	out, err := handlers.ListSystem(app)
	if err != nil {
		fmt.Println("failed to list system variables")
		return
	}

	fmt.Print(string(out))
}

func ListSysPath(app *app.Data) {
	out, err := handlers.ListSysPath(app)
	if err != nil {
		fmt.Println("failed to list path system variable")
		return
	}

	fmt.Print(string(out))
}
