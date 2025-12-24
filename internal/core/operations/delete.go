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

func Delete(payload payload.Data, app *app.Data) {

	ok := inputs.CheckDeleteInput(&payload.Name, &payload.Alias, app)
	if !ok {
		fmt.Println("Cancelling delete...")
		return
	}

	key := payload.Name
	if key == "" {
		if len(payload.Alias) > 0 {
			key = payload.Alias[0]
		}
	}

	command, ok := handlers.DeleteCommand(key, app)
	if !ok {
		app.LoggerInfo.Error("Error delete command: ", payload)
		return
	}

	app.LoggerInfo.Info("deleted command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(notes.GenerateDeleteNote(command))
	history.Save(app)
}
