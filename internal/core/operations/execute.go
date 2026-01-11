package operations

import (
	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func Execute(args []string, app *app.Data) {
	var params []string
	key := ""
	if len(args) > 0 {
		app.Logger.Debug("args", args)
		key = args[0]
		params = args[1:]
	} else {
		ok := inputs.CheckExecuteInput(&key, &params, app)
		if !ok {
			app.Logger.Info("Cancelling execute...")
			return
		}
	}

	output, command, err := handlers.Execute(key, params, app)
	assertExecute(output, command, err, app)
}

func ExecuteLast(op state.MainOperationEnum, args []string, app *app.Data) {
	content := history.GetLast(op, app)

	if content == nil {
		crerrors.PrintCommandNotFoundError(app)
		return
	}

	output, command, err := handlers.Execute(content.CommandName, args, app)
	assertExecute(output, command, err, app)
}

func assertExecute(output string, command *command.Data, err error, app *app.Data) {
	if err != nil {
		app.Logger.Error("Error trying to run command: ", "out", string(output), "err", err)
		return
	}
	app.Logger.Info(string(output))

	app.History.Add(state.Execute, command.Name, notes.GenerateExecuteNote(command))
	history.Save(app)
}
