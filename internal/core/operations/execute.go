package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
)

func Execute(payload data.Payload, app *data.App) {
	output, err := handlers.Execute(payload.Name, payload.Args, app)
	if err != nil {
		app.LoggerInfo.Error("Error trying to run command: ", string(output), err)
		return
	}
	fmt.Println(string(output))
}
