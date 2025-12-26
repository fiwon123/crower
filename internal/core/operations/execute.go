package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func Execute(payload payload.Data, app *app.Data) {
	output, command, err := handlers.Execute(payload.Name, payload.Args, app)
	if err != nil {
		app.LoggerInfo.Error("Error trying to run command: ", string(output), err)
	}
	fmt.Println(string(output))

	app.History.Add(operation.Execute, command.Name, notes.GenerateExecuteNote(command))
	history.Save(app)
}
