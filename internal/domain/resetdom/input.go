package resetdom

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/helper"
	"github.com/fiwon123/crower/internal/interfaces"
)

type Input struct {
	app  *app.Config
	list interfaces.ListHandler
}

func NewInput(app *app.Config, list interfaces.ListHandler) *Input {
	return &Input{
		app:  app,
		list: list,
	}
}

// Verify parameters to process reset operation
func (i *Input) CheckResetInput() bool {
	i.app.Logger.Info("-----------------------------------------")
	i.list.ListCommands()

	i.app.Logger.Info("")
	i.app.Logger.Info("All commands will be erased...")
	ok := helper.GetUserConfirmation("Continue to reset", i.app)
	return ok
}
