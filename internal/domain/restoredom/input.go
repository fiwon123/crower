package restoredom

import (
	"github.com/fiwon123/crower/internal/app"
	historydata "github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/helper"
)

type Input struct {
	app *app.Config
}

func NewInput(app *app.Config) *Input {
	return &Input{
		app: app,
	}
}

// Verify parameters to process restore operation
func (i *Input) CheckRestoreInput() (historydata.Content, bool) {

	var content historydata.Content
	i.app.Logger.Info(i.app.History.GetList())
	content = helper.GetUserInput("Select Row", helper.IsValidContentKey, i.app).(historydata.Content)

	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("Content")
	i.app.Logger.Info("Version:    ", "version", content.Version)
	i.app.Logger.Info("File:    ", "file", content.File)
	i.app.Logger.Info("Note:    ", "note", content.Note)
	i.app.Logger.Info("")

	ok := helper.GetUserConfirmation("Continue to restore", i.app)
	return content, ok
}
