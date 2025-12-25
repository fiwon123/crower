package core

import (
	"fmt"
	"path"

	"github.com/fiwon123/crower/internal/core/operations"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	op "github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"

	"github.com/fiwon123/crower/pkg/utils"
)

// Initialize app based on the cfg file path.
func InitApp(cfgFilePath string) *app.Data {
	var orderKeys []string
	allCommands := command.NewMapData()
	var allAliases command.MapData

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
		allCommands = command.NewMapData()
		allAliases = command.NewMapData()
	}

	return app.New(cfgFilePath, orderKeys, allAliases, allCommands)
}

// Determine which operation will be performed for the user.
func HandlePayload(payload payload.Data, app *app.Data) {
	switch payload.Op {
	case op.Execute:
		operations.Execute(payload, app)
	case op.Create:
		operations.Create(payload, app)
	case op.AddProcess:
		operations.CreateProcess(payload, app)
	case op.Delete:
		operations.Delete(payload, app)
	case op.Update:
		operations.Update(payload, app)
	case op.List:
		operations.List(app)
	case op.Reset:
		operations.Reset(app)
	case op.Open:
		operations.Open(app)
	case op.Process:
		operations.Process(payload, app)
	case op.History:
		operations.History(app)
	case op.Revert:
		operations.Revert(app)
	}
}
