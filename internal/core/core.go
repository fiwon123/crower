package core

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/pkg/utils"
)

// Initialize app based on the cfg file path.
func InitApp(cfgFilePath string) *data.App {
	var allCommands data.CommandsMap
	var allAliases data.CommandsMap

	if cfgFilePath != "" {
		utils.CreateTomlIfNotExists(cfgFilePath)
		fmt.Println("cfgfilepath: ", cfgFilePath)

		err := utils.ReadToml(cfgFilePath, &allCommands)
		if err != nil {
			fmt.Println("error to read toml: ", err)
		}
		allAliases = getAliasMap(allCommands)
	} else {
		allCommands = data.NewCommandsMap()
		allAliases = data.NewCommandsMap()
	}

	return data.NewApp(cfgFilePath, allAliases, allCommands)
}

// Determine which operation will be performed for the user.
func HandlePayload(payload data.Payload, app *data.App) {
	switch payload.Op {
	case data.Execute:
		output, err := handlers.Execute(payload.Name, payload.Args, app)
		if err != nil {
			app.LoggerInfo.Error("Error trying to run command: ", string(output), err)
			return
		}
		fmt.Println(string(output))
	case data.Add:
		err := handlers.AddCommand(payload.Name, payload.Alias, payload.Exec, payload.Args, app)

		if err == nil {
			utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
			app.LoggerInfo.Info("added new command: ", app.AllCommandsByName)
		} else {
			app.LoggerInfo.Error("Error add command: ", err, payload)
		}
	case data.AddProcess:
		err := handlers.AddProcess(payload.Name, payload.Args, app)
		if err == nil {
			utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
			app.LoggerInfo.Info("added new command by process: ", app.AllCommandsByName)
		} else {
			app.LoggerInfo.Error("Error add command by process: ", err, payload)
		}
	case data.Delete:
		if handlers.DeleteCommand(payload.Name, app) {
			app.LoggerInfo.Info("deleted command: ", app.AllCommandsByName)
			utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
		} else {
			app.LoggerInfo.Error("Error delete command: ", payload)
		}
	case data.Update:
		err := handlers.UpdateCommand(payload.Name, payload.Name, payload.Alias, payload.Exec, app)
		if err == nil {
			app.LoggerInfo.Info("updated command: ", app.AllCommandsByName)
			utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
		} else {
			app.LoggerInfo.Error("Error update command: ", err, payload)
		}
	case data.List:
		handlers.List(app)
	case data.Reset:
		handlers.Reset(app)
		utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
		app.LoggerInfo.Info("reset all commands: ", app.AllCommandsByName)
	case data.Open:
		handlers.Open(app.CfgFilePath, app)
	case data.Process:
		handlers.Process(payload.Args, app)
	}
}
