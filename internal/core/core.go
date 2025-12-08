package core

import (
	"fmt"

	"github.com/fiwon123/crower/internal/configuration"
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/pkg/utils"
)

func InitApp(cfgFilePath string) *data.App {
	var commandsMap data.CommandsMap
	var aliasMap data.CommandsMap

	if cfgFilePath != "" {
		utils.CreateTomlIfNotExists(cfgFilePath)
		fmt.Println("cfgfilepath: ", cfgFilePath)

		commandsMap = utils.ReadToml(cfgFilePath)
		aliasMap = getAliasMap(commandsMap)
	} else {
		commandsMap = data.NewCommandsMap()
		aliasMap = data.NewCommandsMap()
	}

	return &data.App{
		CfgFilePath: cfgFilePath,
		AliasMap:    aliasMap,
		CommandsMap: commandsMap,
	}
}

func HandlePayload(payload data.Payload, app *data.App) {
	switch payload.Op {
	case data.Execute:
		input := payload.Command
		command := app.CommandsMap.Get(input.Name)

		if command == nil && len(input.Alias) > 0 {
			fmt.Println("find command by alias ", input.Alias)
			command = app.AliasMap.Get(input.Alias[0])
		}

		if command == nil {
			fmt.Println("command not found")
			return
		}

		output, err := handlers.Execute(*command)
		if err != nil {
			fmt.Println("Error trying to run command: ", err)
			return
		}
		fmt.Println(string(output))
	case data.Add:
		if handlers.AddCommand(payload.Command, app) {
			utils.WriteToml(app.CommandsMap, app.CfgFilePath)
			fmt.Println("added new command: ", app.CommandsMap)
		} else {
			fmt.Println("Error add command: ", payload.Command)
		}
	case data.Delete:
		if handlers.DeleteCommand(payload.Command.Name, app) {
			fmt.Println("deleted command: ", app.CommandsMap)
			utils.WriteToml(app.CommandsMap, app.CfgFilePath)
		} else {
			fmt.Println("Error delete command: ", payload.Command)
		}
	case data.Update:
		if handlers.UpdateCommand(payload.Command, app) {
			fmt.Println("updated command: ", app.CommandsMap)
			utils.WriteToml(app.CommandsMap, app.CfgFilePath)
		} else {
			fmt.Println("Error update command: ", payload.Command)
		}
	case data.List:
		fmt.Println("list all commands: ", app.CommandsMap)
	case data.Reset:
		handlers.Reset(app)
		utils.WriteToml(app.CommandsMap, app.CfgFilePath)
		fmt.Println("reset all commands: ", app.CommandsMap)
	case data.Open:
		configuration.Open(app.CfgFilePath)
	}
}
