package operations

import (
	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func CreateCommand(payload payload.Data, app *app.Data) {
	inputs.CheckCreateInput(&payload.Name, &payload.Alias, &payload.Exec, app)

	command, err := handlers.CreateCommand(payload.Name, payload.Alias, payload.Exec, payload.Args, app)

	if err != nil {
		app.LoggerInfo.Error("Error add command: ", err, payload)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command: ", app.AllCommandsByName)

	app.History.Add(operation.Create, command.Name, notes.GenerateAddNote(command))
	history.Save(app)
}

func CreateProcess(payload payload.Data, app *app.Data) {
	command, err := handlers.CreateProcess(payload.Name, payload.Args, app)
	if err != nil {
		app.LoggerInfo.Error("Error add command by process: ", err, payload)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command by process: ", app.AllCommandsByName)

	app.History.Add(operation.Create, command.Name, notes.GenerateAddProcessNote(command))
	history.Save(app)
}
