package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func Execute(name string, args []string, app *app.Data) {
	output, command, err := handlers.Execute(name, args, app)
	if err != nil {
		app.LoggerInfo.Error("Error trying to run command: ", string(output), err)
	}
	fmt.Println(string(output))

	app.History.Add(state.Execute, command.Name, notes.GenerateExecuteNote(command))
	history.Save(app)
}

func ExecuteLast(op state.OperationEnum, args []string, app *app.Data) {
	content := history.GetLast(op, app)

	if content == nil {
		crerrors.PrintCommandNotFoundError()
		return
	}

	Execute(content.CommandName, args, app)
}
