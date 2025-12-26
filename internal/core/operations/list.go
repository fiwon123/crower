package operations

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/payload"
	"github.com/fiwon123/crower/internal/handlers"
)

func ListCommands(app *app.Data) {
	handlers.ListCommands(app)
}

func ListProcess(payload payload.Data, app *app.Data) {
	handlers.ListProcess(payload.Args, app)
}

func ListHistory(app *app.Data) {
	handlers.ListHistory(app)
}
