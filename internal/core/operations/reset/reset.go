package resetoperations

import (
	resetinputs "github.com/fiwon123/crower/internal/core/inputs/reset"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	resethandlers "github.com/fiwon123/crower/internal/handlers/reset"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

func Reset(app *app.Data) {
	ok := resetinputs.CheckResetInput(app)

	if !ok {
		app.Logger.Info("Cancelling reset...")
		return
	}

	app.Logger.Info("reset all commands: ", app.AllCommandsByName)
	resethandlers.Reset(app)
	crowerutils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(state.Reset, notes.GenerateResetNote())
	history.Save(app)
}
