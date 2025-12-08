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
		command := app.AllCommands.Get(input.Name)

		if command == nil && len(input.AllAlias) > 0 {
			fmt.Println("find command by alias ", input.AllAlias)
			command = app.AllAliases.Get(input.AllAlias[0])
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
			utils.WriteToml(app.AllCommands, app.CfgFilePath)
			fmt.Println("added new command: ", app.AllCommands)
		} else {
			fmt.Println("Error add command: ", payload.Command)
		}
	case data.Delete:
		if handlers.DeleteCommand(payload.Command.Name, app) {
			fmt.Println("deleted command: ", app.AllCommands)
			utils.WriteToml(app.AllCommands, app.CfgFilePath)
		} else {
			fmt.Println("Error delete command: ", payload.Command)
		}
	case data.Update:
		if handlers.UpdateCommand(payload.Command, app) {
			fmt.Println("updated command: ", app.AllCommands)
			utils.WriteToml(app.AllCommands, app.CfgFilePath)
		} else {
			fmt.Println("Error update command: ", payload.Command)
		}
	case data.List:
		fmt.Println("list all commands: ", app.AllCommands)
	case data.Reset:
		handlers.Reset(app)
		utils.WriteToml(app.AllCommands, app.CfgFilePath)
		fmt.Println("reset all commands: ", app.AllCommands)
	case data.Open:
		configuration.Open(app.CfgFilePath)
	}
}
