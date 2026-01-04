package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func ListCommands(app *app.Data) {
	handlers.ListCommands(app)
}

func ListProcess(args []string, app *app.Data) {
	handlers.ListProcess(args, app)
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
	assertListResult(out, err)
}

func ListSystem(app *app.Data) {
	out, err := handlers.ListSystem(app)
	assertListResult(out, err)
}

func ListSysPath(app *app.Data) {
	out, err := handlers.ListSysPath(app)
	assertListResult(out, err)
}

func assertListResult(out []byte, err error) {
	if err != nil {
		fmt.Println("failed to list: ", err, string(out))
		return
	}

	fmt.Print(string(out))
}
