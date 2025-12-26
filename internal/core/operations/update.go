package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func Update(key string, payload payload.Data, app *app.Data) {

	ok := inputs.CheckUpdateInput(&key, &payload.Name, &payload.Alias, &payload.Exec, app)
	if !ok {
		fmt.Println("Cancelling update...")
		return
	}

	oldCommand, newCommand, err := handlers.UpdateCommand(key, payload.Name, payload.Alias, payload.Exec, app)
	if err != nil {
		app.LoggerInfo.Error("Error update command: ", err, payload)
		return
	}

	app.LoggerInfo.Info("updated command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(operation.Update, newCommand.Name, notes.GenerateUpdateNote(oldCommand, newCommand))
	history.Save(app)
}

func UpdateLast(payload payload.Data, app *app.Data) {
	payload.Op = operation.Update
	content := history.GetLast(payload.Op, app)

	if content == nil {
		fmt.Println("command doesn't exist")
		return
	}

	key := content.CommandName

	Update(key, payload, app)
}

func UpdateLastCreate(payload payload.Data, app *app.Data) {
	payload.Op = operation.Create
	content := history.GetLast(payload.Op, app)

	if content == nil {
		fmt.Println("command doesn't exist")
		return
	}

	key := content.CommandName

	payload.Op = operation.Update
	Update(key, payload, app)
}

func UpdateLastExecute(payload payload.Data, app *app.Data) {
	payload.Op = operation.Execute
	content := history.GetLast(payload.Op, app)

	if content == nil {
		fmt.Println("command doesn't exist")
		return
	}

	key := content.CommandName

	payload.Op = operation.Update
	Update(key, payload, app)
}
