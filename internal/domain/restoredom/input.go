package restoredom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/history"
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

// CheckRestoreInput asks which history entry to restore and for confirmation.
//
// It returns false with a nil error when the user declines. An error means the
// choice could not be made at all -- for example an empty history, or no input
// available under --no-input.
func (i *Input) CheckRestoreInput() (history.Content, bool, error) {
	var content history.Content

	i.app.Logger.Info(i.app.History.GetList())

	value, err := helper.GetUserInput("Select Row", helper.IsValidContentKey, i.app)
	if err != nil {
		return content, false, err
	}

	selected, ok := value.(history.Content)
	if !ok {
		return content, false, fmt.Errorf("selected row is not a history entry")
	}
	content = selected

	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("Content")
	i.app.Logger.Info("Version:    ", content.Version)
	i.app.Logger.Info("File:    ", content.File)
	i.app.Logger.Info("Note:    ", content.Note)
	i.app.Logger.Info("")

	ok, confirmErr := helper.GetUserConfirmation("Continue to restore", i.app)
	return content, ok, confirmErr
}
