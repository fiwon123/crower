package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/payload"
	"github.com/fiwon123/crower/internal/handlers"
)

func Execute(payload payload.Data, app *app.Data) {
	output, err := handlers.Execute(payload.Name, payload.Args, app)
	if err != nil {
		app.LoggerInfo.Error("Error trying to run command: ", string(output), err)
		return
	}
	fmt.Println(string(output))
}
