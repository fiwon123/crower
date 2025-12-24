package operations

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func List(app *app.Data) {
	handlers.List(app)
}
