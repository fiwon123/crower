package operations

import (
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
)

func Process(payload data.Payload, app *data.App) {
	handlers.Process(payload.Args, app)
}
