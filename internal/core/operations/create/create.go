package createoperations

import (
	"path/filepath"
	"runtime"
	"strings"

	createinputs "github.com/fiwon123/crower/internal/core/inputs/create"
	"github.com/fiwon123/crower/internal/crowererrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/state"
	createhandlers "github.com/fiwon123/crower/internal/handlers/create"
	openhandlers "github.com/fiwon123/crower/internal/handlers/open"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

func CreateCommand(allAlias []string, args []string, app *app.Data) {
	name := ""
	exec := ""
	if len(args) == 2 {
		name = args[0]
		exec = args[1]
	} else {
		createinputs.CheckCreateInput(&name, &allAlias, &exec, app)
	}

	command := performCreateCommand(name, allAlias, exec, app)
	if command == nil {
		return
	}

	app.History.Add(state.Create, notes.GenerateCreateCommandNote(command, args))
	history.Save(app)
}

func performCreateCommand(name string, allAlias []string, exec string, app *app.Data) *command.Data {
	command, err := createhandlers.CreateCommand(name, allAlias, exec, app)

	if err != nil {
		app.Logger.Error("Error add command: ", "error", err, "name", name, "alias", allAlias, "exec", exec)
		return nil
	}

	crowerutils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.Logger.Info("added new command: ", "allCommands", app.AllCommandsByName)

	return command
}

func CreateProcess(name string, args []string, app *app.Data) {
	command, err := createhandlers.CreateProcess(name, args, app)
	if err != nil {
		app.Logger.Error("Error add command by process: ", "error", err, "name", name, "args", args)
		return
	}

	crowerutils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
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
		crowererrors.PrintNotArgs("var name and var value", app)
		return
	}

	out, err := createhandlers.CreateSystemVariable(newVar, value, app)
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
		crowererrors.PrintNotArgs("path", app)
		return
	}

	out, err := createhandlers.CreateSystemPathVariable(newPath, app)
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
		err := createhandlers.CreateFile(path, app)
		if err != nil {
			app.Logger.Error(err.Error())
		}
	}

	app.History.Add(state.Create, notes.GenerateCreateFile(args))
	history.Save(app)
}

func CreateFolder(args []string, app *app.Data) {
	for _, path := range args {
		err := createhandlers.CreateFolder(path, app)
		if err != nil {
			app.Logger.Error(err.Error())
		}
	}

	app.History.Add(state.Create, notes.GenerateCreateFolder(args))
	history.Save(app)
}

func CreateLastCommand(op state.MainOperationEnum, args []string, app *app.Data) {
	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		crowererrors.PrintNotArgs("name", app)
		return
	}

	content := history.GetLast(op, app)

	if content == nil {
		crowererrors.PrintCommandNotFoundError(app)
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

func CreateScriptCommand(args []string, app *app.Data) {
	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		crowererrors.PrintNotArgs("name", app)
	}

	scriptFilePath, err := createhandlers.CreateScriptCommand(name, app)
	if err != nil {
		app.Logger.Error(err.Error())
		return
	}

	var command *command.Data
	switch runtime.GOOS {
	case "windows":
		command = performCreateCommand(name, []string{}, scriptFilePath, app)
	case "linux":
		command = performCreateCommand(name, []string{}, scriptFilePath, app)
	}

	if command == nil {
		return
	}

	openhandlers.Open([]string{filepath.Dir(scriptFilePath)}, app)

	app.History.Add(state.Create, notes.GenerateCreateScriptCommandNote(command, args))
	history.Save(app)
}
