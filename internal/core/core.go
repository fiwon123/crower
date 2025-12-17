package core

import (
	"fmt"

	"github.com/fiwon123/crower/internal/configuration"
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/pkg/utils"
)

func InitApp(cfgFilePath string) *data.App {
	var allCommands data.CommandsMap
	var allAliases data.CommandsMap

	if cfgFilePath != "" {
		utils.CreateTomlIfNotExists(cfgFilePath)
		fmt.Println("cfgfilepath: ", cfgFilePath)

		allCommands = utils.ReadToml(cfgFilePath)
		allAliases = getAliasMap(allCommands)
	} else {
		allCommands = data.NewCommandsMap()
		allAliases = data.NewCommandsMap()
	}

	return data.NewApp(cfgFilePath, allAliases, allCommands)
}

func HandlePayload(payload data.Payload, app *data.App) {
	switch payload.Op {
	case data.Execute:
		input := payload.Command
		command := app.AllCommandsByName.Get(input.Name)

		if command == nil && len(input.AllAlias) > 0 {
			app.LoggerInfo.Info("find command by alias ", input.AllAlias)
			command = app.AllCommandsByAlias.Get(input.AllAlias[0])
		}

		if command == nil {
			app.LoggerInfo.Error("command not found")
			return
		}

		output, err := handlers.Execute(*command, app)
		if err != nil {
			app.LoggerInfo.Error("Error trying to run command: ", err)
			return
		}
		fmt.Println(string(output))
	case data.Add:
		if handlers.AddCommand(payload.Command, app) {
			utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
			app.LoggerInfo.Info("added new command: ", app.AllCommandsByName)
		} else {
			app.LoggerInfo.Error("Error add command: ", payload.Command)
		}
	case data.Delete:
		if handlers.DeleteCommand(payload.Command.Name, app) {
			app.LoggerInfo.Info("deleted command: ", app.AllCommandsByName)
			utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
		} else {
			app.LoggerInfo.Error("Error delete command: ", payload.Command)
		}
	case data.Update:
		if handlers.UpdateCommand(payload.Command.Name, payload.Command, app) {
			app.LoggerInfo.Info("updated command: ", app.AllCommandsByName)
			utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
		} else {
			app.LoggerInfo.Error("Error update command: ", payload.Command)
		}
	case data.List:
		app.LoggerInfo.Info("list all commands: ", app.AllCommandsByName)
	case data.Reset:
		handlers.Reset(app)
		utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
		app.LoggerInfo.Info("reset all commands: ", app.AllCommandsByName)
	case data.Open:
		configuration.Open(app.CfgFilePath, app)
	}
}
