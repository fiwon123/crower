package resetinputs

import (
	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

// Verify parameters to process reset operation
func CheckResetInput(app *app.Data) bool {
	app.Logger.Info("-----------------------------------------")
	handlers.ListCommands(app)

	app.Logger.Info("")
	app.Logger.Info("All commands will be erased...")
	ok := inputs.GetUserConfirmation("Continue to reset", app)
	return ok
}
