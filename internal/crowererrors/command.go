package crowererrors

import (
	"github.com/fiwon123/crower/internal/data/app"
)

func PrintCommandNotFoundError(app *app.Data) {
	app.Logger.Error("command not found")
}
