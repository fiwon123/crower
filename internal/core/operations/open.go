package operations

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func Open(app *app.Data) {
	handlers.Open(app.CfgFilePath, app)
}
