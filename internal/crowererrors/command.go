package crowererrors

import (
	appdata "github.com/fiwon123/crower/internal/data/app"
)

func PrintCommandNotFoundError(app *appdata.Data) {
	app.Logger.Error("command not found")
}
