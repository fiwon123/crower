package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"

	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func Delete(args []string, app *app.Data) {

	key := ""
	if len(args) > 0 {
		key = args[0]
	}

	ok := inputs.CheckDeleteInput(&key, app)
	if !ok {
		fmt.Println("Cancelling delete...")
		return
	}

	command, ok := handlers.DeleteCommand(key, app)
	if !ok {
		app.LoggerInfo.Error("Error delete command: ", key)
		return
	}

	app.LoggerInfo.Info("deleted command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(state.Delete, command.Name, notes.GenerateDeleteNote(command))
	history.Save(app)
}

func DeleteLast(op state.OperationEnum, app *app.Data) {
	content := history.GetLast(op, app)

	if content == nil {
		crerrors.PrintCommandNotFoundError()
		return
	}

	Delete([]string{content.CommandName}, app)
}

func DeleteSystemVariable(args []string, app *app.Data) {
	newVar := ""
	if len(args) >= 1 {
		newVar = args[0]
	} else {
		crerrors.PrintNotArgs("var name")
		return
	}

	out, err := handlers.DeleteSystemVariable(newVar, app)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(out)
}

func DeleteSystemPathVariable(args []string, app *app.Data) {
	newPath := ""
	if len(args) > 0 {
		newPath = args[0]
	} else {
		crerrors.PrintNotArgs("path")
		return
	}

	out, err := handlers.DeleteSystemPathVariable(newPath, app)
	if err != nil {
		fmt.Printf("err: %s \n", err)
		return
	}

	fmt.Println(out)
}

func DeleteFile(args []string, app *app.Data) {
	filePath := ""
	if len(args) > 0 {
		filePath = args[0]
	} else {
		crerrors.PrintNotFileAndOutputPath()
		return
	}

	handlers.DeleteFile(filePath, app)
}

func DeleteFolder(args []string, app *app.Data) {
	folderPath := ""
	if len(args) > 0 {
		folderPath = args[0]
	} else {
		crerrors.PrintNotFileAndOutputPath()
		return
	}

	handlers.DeleteFolder(folderPath, app)
}
