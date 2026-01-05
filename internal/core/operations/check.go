package operations

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func CheckNewVersion(currentVersion string, app *app.Data) {
	handlers.CheckNewVersion(currentVersion, app)
}
