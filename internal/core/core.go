package core

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/pkg/utils"
)

func HandlePayload(payload data.Payload, app *data.App) {
	switch payload.Op {
	case data.Execute:
		output, err := handlers.Execute(payload.Command, app)
		if err != nil {
			fmt.Println("Error trying to run command: ", err)
			return
		}
		fmt.Println(string(output))
	case data.Create:
		handlers.AddCommand(payload.Command, app)
		utils.WriteToml(app.CommandsMap, app.CfgFilePath)
		fmt.Println("added new command: ", app.CommandsMap)
	case data.Delete:
	case data.Update:
	case data.List:
		fmt.Println("list all commands: ", app.CommandsMap)
	case data.Reset:
		handlers.Reset(app)
		utils.WriteToml(app.CommandsMap, app.CfgFilePath)
		fmt.Println("reset all commands: ", app.CommandsMap)
	}
}
