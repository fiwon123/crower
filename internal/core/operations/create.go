package operations

import (
	"strings"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
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

	command := performCreateCommand(name, allAlias, exec, app)
	if command == nil {
		return
	}

	app.History.Add(state.Create, notes.GenerateCreateCommandNote(command, args))
	history.Save(app)
}

func performCreateCommand(name string, allAlias []string, exec string, app *app.Data) *command.Data {
	command, err := handlers.CreateCommand(name, allAlias, exec, app)

	if err != nil {
		app.Logger.Error("Error add command: ", "error", err, "name", name, "alias", allAlias, "exec", exec)
		return nil
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.Logger.Info("added new command: ", app.AllCommandsByName)

	return command
}

func CreateProcess(name string, args []string, app *app.Data) {
	command, err := handlers.CreateProcess(name, args, app)
	if err != nil {
		app.Logger.Error("Error add command by process: ", "error", err, "name", name, "args", args)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.Logger.Info("added new command by process: ", "allCommands", app.AllCommandsByName)

	app.History.Add(state.Create, notes.GenerateCreateProcessNote(command, args))
	history.Save(app)
}

func CreateSystemVariable(args []string, app *app.Data) {
	newVar := ""
	value := ""
	if len(args) >= 2 {
		newVar = args[0]
		value = args[1]
	} else {
		crerrors.PrintNotArgs("var name and var value", app)
		return
	}

	out, err := handlers.CreateSystemVariable(newVar, value, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	app.Logger.Info(out)

	app.History.Add(state.Create, notes.GenerateCreateSystemVariableNote(args))
	history.Save(app)
}

func CreateSystemPathVariable(args []string, app *app.Data) {
	newPath := ""
	if len(args) > 0 {
		newPath = args[0]
	} else {
		crerrors.PrintNotArgs("path", app)
		return
	}

	out, err := handlers.CreateSystemPathVariable(newPath, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	app.Logger.Info(out)

	app.History.Add(state.Create, notes.GenerateCreateSystemPathVariableNote(args))
	history.Save(app)
}

func CreateFile(args []string, app *app.Data) {
	for _, path := range args {
		err := handlers.CreateFile(path, app)
		if err != nil {
			app.Logger.Error(err.Error())
		}
	}

	app.History.Add(state.Create, notes.GenerateCreateFile(args))
	history.Save(app)
}

func CreateFolder(args []string, app *app.Data) {
	for _, path := range args {
		err := handlers.CreateFolder(path, app)
		if err != nil {
			app.Logger.Error(err.Error())
		}
	}

	app.History.Add(state.Create, notes.GenerateCreateFolder(args))
	history.Save(app)
}

func CreateLastCommand(op state.MainOperationEnum, name string, app *app.Data) {
	content := history.GetLast(op, app)

	if content == nil {
		crerrors.PrintCommandNotFoundError(app)
		return
	}

	exec := ""
	key := content.CommandName
	if key == "" {
		splitted := strings.SplitSeq(content.Note, ";")
		for keyValRaw := range splitted {
			keyVal := strings.Split(keyValRaw, "=")
			if keyVal[0] == "exec" {
				exec = keyVal[1]
				break
			}
		}
	}

	command := performCreateCommand(name, []string{}, exec, app)
	if command == nil {
		return
	}

	app.History.Add(state.Create, notes.GenerateCreateCommandLastExecuteNote(command))
	history.Save(app)
}
