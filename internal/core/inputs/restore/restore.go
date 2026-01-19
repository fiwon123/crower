package restoreinputs

import (
	inputscore "github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/history"
)

// Verify parameters to process restore operation
func CheckRestoreInput(app *app.Data) (history.Content, bool) {

	var content history.Content
	app.Logger.Info(app.History.GetList())
	content = inputscore.GetUserInput("Select Row", inputscore.IsValidContentKey, app).(history.Content)

	app.Logger.Info("-----------------------------------------")
	app.Logger.Info("Content")
	app.Logger.Info("Version:    ", "version", content.Version)
	app.Logger.Info("File:    ", "file", content.File)
	app.Logger.Info("Note:    ", "note", content.Note)
	app.Logger.Info("")

	ok := inputscore.GetUserConfirmation("Continue to restore", app)
	return content, ok
}
