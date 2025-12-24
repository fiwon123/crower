package operations

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func Reset(app *app.Data) {
	app.LoggerInfo.Info("reset all commands: ", app.AllCommandsByName)
	handlers.Reset(app)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(notes.GenerateResetNote())
	history.Save(app)
}
