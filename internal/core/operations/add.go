package operations

import (
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/pkg/utils"
)

func AddOp(payload data.Payload, app *data.App) {
	command, err := handlers.AddCommand(payload.Name, payload.Alias, payload.Exec, payload.Args, app)

	if err != nil {
		app.LoggerInfo.Error("Error add command: ", err, payload)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command: ", app.AllCommandsByName)

	app.History.Add(history.GenerateAddNote(command))
	history.Save(app)
}
