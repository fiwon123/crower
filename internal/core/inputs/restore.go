package inputs

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/history"
)

// Verify parameters to process restore operation
func CheckRestoreInput(app *app.Data) (history.Content, bool) {

	var content history.Content
	app.Logger.Info(app.History.GetList())
	content = getUserInput("Select Row", isValidContentKey, app).(history.Content)

	app.Logger.Info("-----------------------------------------")
	app.Logger.Info("Content")
	app.Logger.Info("Version:    ", "version", content.Version)
	app.Logger.Info("File:    ", "file", content.File)
	app.Logger.Info("Note:    ", "note", content.Note)
	app.Logger.Info("")

	ok := getUserConfirmation("Continue to restore", app)
	return content, ok
}
