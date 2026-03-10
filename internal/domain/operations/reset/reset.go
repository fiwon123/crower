package resetoperations

import (
	resetinputs "github.com/fiwon123/crower/internal/core/inputs/reset"
	appdata "github.com/fiwon123/crower/internal/data/app"
	resetnotesdata "github.com/fiwon123/crower/internal/data/notes/reset"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	resethandlers "github.com/fiwon123/crower/internal/handlers/reset"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

func Reset(app *appdata.Data) {
	ok := resetinputs.CheckResetInput(app)

	if !ok {
		app.Logger.Info("Cancelling reset...")
		return
	}

	app.Logger.Info("reset all commands: ", app.AllCommandsByName)
	resethandlers.Reset(app)
	crowerutils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(operationsdata.Reset, resetnotesdata.NewResetNote())
	historyhelper.Save(app)
}
