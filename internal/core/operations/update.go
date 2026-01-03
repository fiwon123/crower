package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/cterrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func Update(key string, name string, allAlias []string, exec string, app *app.Data) {

	ok := inputs.CheckUpdateInput(&key, &name, &allAlias, &exec, app)
	if !ok {
		fmt.Println("Cancelling update...")
		return
	}

	oldCommand, newCommand, err := handlers.UpdateCommand(key, name, allAlias, exec, app)
	if err != nil {
		app.LoggerInfo.Error("Error update command: ", err, key, name, allAlias, exec)
		return
	}

	app.LoggerInfo.Info("updated command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(operation.Update, newCommand.Name, notes.GenerateUpdateNote(oldCommand, newCommand))
	history.Save(app)
}

func UpdateLast(op operation.State, name string, allAlias []string, exec string, app *app.Data) {
	content := history.GetLast(op, app)

	if content == nil {
		cterrors.PrintCommandNotFoundError()
		return
	}

	key := content.CommandName

	Update(key, name, allAlias, exec, app)
}
