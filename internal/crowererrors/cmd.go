package crowererrors

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
)

func PrintCmdHelp(cmdName string, app *app.Data) {
	app.Logger.Info(fmt.Sprintf("Type 'crower %s --help' for more information", cmdName))
}

func PrintEmptyArgs(app *app.Data) {
	app.Logger.Info("args is empty")
}
