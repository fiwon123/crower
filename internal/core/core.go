package core

import (
	"fmt"
	"path"

	"github.com/fiwon123/crower/internal/core/operations"
	"github.com/fiwon123/crower/internal/data"
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
		operations.Execute(payload, app)
	case data.AddOp:
		operations.AddOp(payload, app)
	case data.AddProcess:
		operations.AddProcess(payload, app)
	case data.DeleteOp:
		operations.Delete(payload, app)
	case data.UpdateOp:
		operations.Update(payload, app)
	case data.ListOp:
		operations.List(app)
	case data.ResetOp:
		operations.Reset(app)
	case data.OpenOp:
		operations.Open(app)
	case data.ProcessOp:
		operations.Process(payload, app)
	case data.HistoryOp:
		operations.History(app)
	case data.RevertOp:
		operations.Revert(app)
	}
}
