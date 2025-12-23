package operations

import (
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/pkg/utils"
)

func AddProcess(payload data.Payload, app *data.App) {
	command, err := handlers.AddProcess(payload.Name, payload.Args, app)
	if err != nil {
		app.LoggerInfo.Error("Error add command by process: ", err, payload)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command by process: ", app.AllCommandsByName)

	app.History.Add(history.GenerateAddProcessNote(command))
	history.Save(app)
}
