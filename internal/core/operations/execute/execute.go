package executeoperations

import (
	executeinputs "github.com/fiwon123/crower/internal/core/inputs/execute"
	"github.com/fiwon123/crower/internal/crowererrors"
	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	executenotesdata "github.com/fiwon123/crower/internal/data/notes/execute"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	executehandlers "github.com/fiwon123/crower/internal/handlers/execute"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
)

func ExecuteCommand(args []string, app *appdata.Data) {
	var params []string
	key := ""
	if len(args) > 0 {
		app.Logger.Debug("args", args)
		key = args[0]
		params = args[1:]
	} else {
		ok := executeinputs.CheckExecuteInput(&key, &params, app)
		if !ok {
			app.Logger.Info("Cancelling execute...")
			return
		}
	}

	output, command, err := executehandlers.Execute(key, params, app)
	assertExecute(output, command, err, app)

	app.History.Add(operationsdata.Execute, executenotesdata.NewExecuteCommandNote(command))
	historyhelper.Save(app)
}

func ExecuteLast(op operationsdata.MainOperationEnum, args []string, app *appdata.Data) {
	content := historyhelper.GetLast(op, app)

	if content == nil {
		crowererrors.PrintCommandNotFoundError(app)
		return
	}

	output, command, err := executehandlers.Execute(content.CommandName, args, app)
	assertExecute(output, command, err, app)

	app.History.Add(operationsdata.Execute, executenotesdata.NewExecuteLastNote(op, command))
	historyhelper.Save(app)
}

func assertExecute(output string, command *commanddata.Data, err error, app *appdata.Data) {
	if err != nil {
		app.Logger.Error("Error trying to run command: ", "out", string(output), "err", err)
		return
	}
	app.Logger.Info(string(output))
}
