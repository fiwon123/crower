package resetinputs

import (
	inputscore "github.com/fiwon123/crower/internal/core/inputs"
	appdata "github.com/fiwon123/crower/internal/data/app"
	listhandlers "github.com/fiwon123/crower/internal/handlers/list"
)

// Verify parameters to process reset operation
func CheckResetInput(app *appdata.Data) bool {
	app.Logger.Info("-----------------------------------------")
	listhandlers.ListCommands(app)

	app.Logger.Info("")
	app.Logger.Info("All commands will be erased...")
	ok := inputscore.GetUserConfirmation("Continue to reset", app)
	return ok
}
