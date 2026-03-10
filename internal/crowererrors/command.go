package crowererrors

import "github.com/fiwon123/crower/internal/app"

func PrintCommandNotFoundError(app *app.Config) {
	app.Logger.Error("command not found")
}
