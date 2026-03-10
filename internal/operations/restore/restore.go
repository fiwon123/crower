package restoreoperations

import (
	restoreinputs "github.com/fiwon123/crower/internal/core/inputs/restore"
	appdata "github.com/fiwon123/crower/internal/data/app"
	restorenotesdata "github.com/fiwon123/crower/internal/data/notes/restore"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	restorehandlers "github.com/fiwon123/crower/internal/handlers/restore"

	historyhelper "github.com/fiwon123/crower/internal/helper/history"
)

func Restore(args []string, app *appdata.Data) {
	key := ""
	if len(args) > 0 {
		key = args[0]
	}

	content, ok := restoreinputs.CheckRestoreInput(app)
	if !ok {
		app.Logger.Info("Cancelling Restore...")
		return
	}

	out, err := restorehandlers.RestoreHistory(key, content, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	app.Logger.Info("restored command: ", "out", out)

	app.History.Add(operationsdata.Restore, restorenotesdata.NewRestoreNote(out))
	historyhelper.Save(app)
}
