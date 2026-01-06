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

func CreateCommand(allAlias []string, args []string, app *app.Data) {
	name := ""
	exec := ""
	if len(args) == 2 {
		name = args[0]
		exec = args[1]
	} else {
		inputs.CheckCreateInput(&name, &allAlias, &exec, app)
	}

	command, err := handlers.CreateCommand(name, allAlias, exec, app)

	if err != nil {
		app.LoggerInfo.Error("Error add command: ", err, name, allAlias, exec, args)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command: ", app.AllCommandsByName)

	app.History.Add(state.Create, command.Name, notes.GenerateAddNote(command))
	history.Save(app)
}

func CreateProcess(name string, args []string, app *app.Data) {
	command, err := handlers.CreateProcess(name, args, app)
	if err != nil {
		app.LoggerInfo.Error("Error add command by process: ", err, name, args)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command by process: ", app.AllCommandsByName)

	app.History.Add(state.Create, command.Name, notes.GenerateAddProcessNote(command))
	history.Save(app)
}

func CreateSystemVariable(args []string, app *app.Data) {
	newVar := ""
	value := ""
	if len(args) >= 2 {
		newVar = args[0]
		value = args[1]
	} else {
		crerrors.PrintNotArgs("var name and var value")
		return
	}

	out, err := handlers.CreateSystemVariable(newVar, value, app)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(out)
}

func CreateSystemPathVariable(args []string, app *app.Data) {
	newPath := ""
	if len(args) > 0 {
		newPath = args[0]
	} else {
		crerrors.PrintNotArgs("path")
		return
	}

	out, err := handlers.CreateSystemPathVariable(newPath, app)
	if err != nil {
		fmt.Printf("err: %s \n", err)
		return
	}

	fmt.Println(out)
}

func CreateFile(args []string, app *app.Data) {
	for _, path := range args {
		err := handlers.CreateFile(path, app)
		if err != nil {
			fmt.Printf("err: %s \n", err)
		}
	}
}

func CreateFolder(args []string, app *app.Data) {
	for _, path := range args {
		err := handlers.CreateFolder(path, app)
		if err != nil {
			fmt.Println(err)
		}
	}
}
