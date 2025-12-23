package operations

import (
	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/pkg/utils"
)

func Update(payload data.Payload, app *data.App) {
	key := ""
	if len(payload.Args) != 0 {
		key = payload.Args[0]
	}
	err := inputs.CheckUpdateInput(&key, &payload.Name, &payload.Alias, &payload.Exec, app)
	if err != nil {
		app.LoggerInfo.Error("Error update command: ", err, payload)
		return
	}

	oldCommand, newCommand, err := handlers.UpdateCommand(key, payload.Name, payload.Alias, payload.Exec, app)
	if err != nil {
		app.LoggerInfo.Error("Error update command: ", err, payload)
		return
	}

	app.LoggerInfo.Info("updated command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(history.GenerateUpdateNote(oldCommand, newCommand))
	history.Save(app)
}
