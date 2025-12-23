package operations

import (
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/pkg/utils"
)

func Delete(payload data.Payload, app *data.App) {
	command, ok := handlers.DeleteCommand(payload.Name, app)
	if !ok {
		app.LoggerInfo.Error("Error delete command: ", payload)
		return
	}

	app.LoggerInfo.Info("deleted command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(history.GenerateDeleteNote(command))
	history.Save(app)
}
