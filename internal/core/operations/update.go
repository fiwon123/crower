package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/payload"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func Update(payload payload.Data, app *app.Data) {
	key := ""
	if len(payload.Args) != 0 {
		key = payload.Args[0]
	}

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

	app.History.Add(notes.GenerateUpdateNote(oldCommand, newCommand))
	history.Save(app)
}
