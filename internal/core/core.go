package core

import (
	"fmt"
	"path"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/pkg/utils"
)

// Initialize app based on the cfg file path.
func InitApp(cfgFilePath string) *data.App {
	var orderKeys []string
	allCommands := data.NewCommandsMap()
	var allAliases data.CommandsMap

	if cfgFilePath != "" {
		utils.CreateFolderIfNotExists(path.Dir(cfgFilePath))
		utils.CreateFileIfNotExists(cfgFilePath)
		fmt.Println("cfgfilepath: ", cfgFilePath)

		var err error
		orderKeys, err = utils.ReadKeysTomlInOrder(cfgFilePath)
		err = utils.ReadToml(cfgFilePath, &allCommands)
		if err != nil {
			fmt.Println("error to read toml: ", err)
		}
		allAliases = getAliasMap(allCommands)
	} else {
		allCommands = data.NewCommandsMap()
		allAliases = data.NewCommandsMap()
	}

	return data.NewApp(cfgFilePath, orderKeys, allAliases, allCommands)
}

// Determine which operation will be performed for the user.
func HandlePayload(payload data.Payload, app *data.App) {
	switch payload.Op {
	case data.ExecuteOp:
		executeOp(payload, app)
	case data.AddOp:
		addOp(payload, app)
	case data.AddProcess:
		addProcess(payload, app)
	case data.DeleteOp:
		deleteOp(payload, app)
	case data.UpdateOp:
		updateOp(payload, app)
	case data.ListOp:
		listOp(app)
	case data.ResetOp:
		resetOp(app)
	case data.OpenOp:
		openOp(app)
	case data.ProcessOp:
		processOp(payload, app)
	case data.HistoryOp:
		historyOp(app)
	}
}

func historyOp(app *data.App) {
	handlers.History(app)
}

func processOp(payload data.Payload, app *data.App) {
	handlers.Process(payload.Args, app)
}

func openOp(app *data.App) {
	handlers.Open(app.CfgFilePath, app)
}

func resetOp(app *data.App) {
	handlers.Reset(app)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add("Reset")
	handlers.SaveHistory(app)

	app.LoggerInfo.Info("reset all commands: ", app.AllCommandsByName)
}

func listOp(app *data.App) {
	handlers.List(app)
}

func updateOp(payload data.Payload, app *data.App) {
	key := ""
	if len(payload.Args) != 0 {
		key = payload.Args[0]
	}
	err := checkInputUpdate(&key, &payload.Name, &payload.Alias, &payload.Exec, app)
	if err != nil {
		app.LoggerInfo.Error("Error update command: ", err, payload)
	}

	err = handlers.UpdateCommand(key, payload.Name, payload.Alias, payload.Exec, app)
	if err == nil {
		app.LoggerInfo.Info("updated command: ", app.AllCommandsByName)
		utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

		app.History.Add("Update")
		handlers.SaveHistory(app)
	} else {
		app.LoggerInfo.Error("Error update command: ", err, payload)
	}
}

func deleteOp(payload data.Payload, app *data.App) {
	if handlers.DeleteCommand(payload.Name, app) {
		app.LoggerInfo.Info("deleted command: ", app.AllCommandsByName)
		utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

		app.History.Add("Delete")
		handlers.SaveHistory(app)
	} else {
		app.LoggerInfo.Error("Error delete command: ", payload)
	}
}

func addProcess(payload data.Payload, app *data.App) {
	err := handlers.AddProcess(payload.Name, payload.Args, app)
	if err == nil {
		utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
		app.LoggerInfo.Info("added new command by process: ", app.AllCommandsByName)

		app.History.Add("Add Process")
		handlers.SaveHistory(app)
	} else {
		app.LoggerInfo.Error("Error add command by process: ", err, payload)
	}
}

func addOp(payload data.Payload, app *data.App) {
	err := handlers.AddCommand(payload.Name, payload.Alias, payload.Exec, payload.Args, app)

	if err == nil {
		utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
		app.LoggerInfo.Info("added new command: ", app.AllCommandsByName)

		app.History.Add("Add")
		handlers.SaveHistory(app)
	} else {
		app.LoggerInfo.Error("Error add command: ", err, payload)
	}
}

func executeOp(payload data.Payload, app *data.App) {
	output, err := handlers.Execute(payload.Name, payload.Args, app)
	if err != nil {
		app.LoggerInfo.Error("Error trying to run command: ", string(output), err)
		return
	}
	fmt.Println(string(output))
}
