package core

import (
	"fmt"
	"path"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/inputs"
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
	case data.RevertOp:
		revertOp(app)
	}
}

func revertOp(app *data.App) {
	backHistory := app.History.GetBeforeLast()

	if backHistory == nil {
		app.LoggerInfo.Error("Error already in initial history")
		return
	}

	err := history.RevertTo(backHistory, app)
	if err != nil {
		app.LoggerInfo.Error("Error revert history %v", err)
		return
	}
	app.LoggerInfo.Info("reverted to history version ", backHistory.Version)
	history.Save(app)
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
	app.LoggerInfo.Info("reset all commands: ", app.AllCommandsByName)
	handlers.Reset(app)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(history.GenerateResetNote())
	history.Save(app)
}

func listOp(app *data.App) {
	handlers.List(app)
}

func updateOp(payload data.Payload, app *data.App) {
	key := ""
	if len(payload.Args) != 0 {
		key = payload.Args[0]
	}
	err := inputs.CheckUpdateInput(&key, &payload.Name, &payload.Alias, &payload.Exec, app)
	if err != nil {
		app.LoggerInfo.Error("Error update command: ", err, payload)
		return
	}

	oldCommand, newCommand, err := handlers.UpdateCommand(key, payload.Name, payload.Alias, payload.Exec, app)
	if err != nil {
		app.LoggerInfo.Error("Error update command: ", err, payload)
		return
	}

	app.LoggerInfo.Info("updated command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(history.GenerateUpdateNote(oldCommand, newCommand))
	history.Save(app)
}

func deleteOp(payload data.Payload, app *data.App) {
	command, ok := handlers.DeleteCommand(payload.Name, app)
	if !ok {
		app.LoggerInfo.Error("Error delete command: ", payload)
		return
	}

	app.LoggerInfo.Info("deleted command: ", app.AllCommandsByName)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(history.GenerateDeleteNote(command))
	history.Save(app)
}

func addProcess(payload data.Payload, app *data.App) {
	command, err := handlers.AddProcess(payload.Name, payload.Args, app)
	if err != nil {
		app.LoggerInfo.Error("Error add command by process: ", err, payload)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command by process: ", app.AllCommandsByName)

	app.History.Add(history.GenerateAddProcessNote(command))
	history.Save(app)
}

func addOp(payload data.Payload, app *data.App) {
	command, err := handlers.AddCommand(payload.Name, payload.Alias, payload.Exec, payload.Args, app)

	if err != nil {
		app.LoggerInfo.Error("Error add command: ", err, payload)
		return
	}

	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)
	app.LoggerInfo.Info("added new command: ", app.AllCommandsByName)

	app.History.Add(history.GenerateAddNote(command))
	history.Save(app)
}

func executeOp(payload data.Payload, app *data.App) {
	output, err := handlers.Execute(payload.Name, payload.Args, app)
	if err != nil {
		app.LoggerInfo.Error("Error trying to run command: ", string(output), err)
		return
	}
	fmt.Println(string(output))
}
