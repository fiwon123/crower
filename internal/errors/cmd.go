package errors

import (
	"fmt"

	"github.com/fiwon123/crower/internal/app"
)

func PrintCmdHelp(cmdName string, app *app.Config) {
	app.Logger.Info(fmt.Sprintf("Type 'crower %s --help' for more information", cmdName))
}

func PrintEmptyArgs(app *app.Config) {
	app.Logger.Info("args is empty")
}
