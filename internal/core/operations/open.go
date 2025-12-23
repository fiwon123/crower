package operations

import (
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
)

func Open(app *data.App) {
	handlers.Open(app.CfgFilePath, app)
}
