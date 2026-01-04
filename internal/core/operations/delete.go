package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/cterrors"
	"github.com/fiwon123/crower/internal/data/app"

	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func Delete(name string, allAlias []string, app *app.Data) {

	ok := inputs.CheckDeleteInput(&name, &allAlias, app)
	if !ok {
		fmt.Println("Cancelling delete...")
		return
	}

	key := name
	if key == "" {
		if len(allAlias) > 0 {
			key = allAlias[0]
		}
	}

	command, ok := handlers.DeleteCommand(key, app)
	if !ok {
		app.LoggerInfo.Error("Error delete command: ", name, allAlias)
		return
	}

	app.LoggerInfo.Info("deleted command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(state.Delete, command.Name, notes.GenerateDeleteNote(command))
	history.Save(app)
}

func DeleteLast(op state.OperationEnum, name string, allAlias []string, app *app.Data) {
	content := history.GetLast(op, app)

	if content == nil {
		cterrors.PrintCommandNotFoundError()
		return
	}

	Delete(name, allAlias, app)
}

func DeleteFile(args []string, app *app.Data) {
	filePath := ""
	if len(args) > 0 {
		filePath = args[0]
	} else {
		cterrors.PrintNotFileAndOutputPath()
		return
	}

	handlers.DeleteFile(filePath, app)
}

func DeleteFolder(args []string, app *app.Data) {
	folderPath := ""
	if len(args) > 0 {
		folderPath = args[0]
	} else {
		cterrors.PrintNotFileAndOutputPath()
		return
	}

	handlers.DeleteFolder(folderPath, app)
}
