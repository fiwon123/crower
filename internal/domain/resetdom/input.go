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

// CheckResetInput shows what will be erased and asks for confirmation.
func (i *Input) CheckResetInput() (bool, error) {
	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info(i.list.CommandsTable())

	i.app.Logger.Info("")
	i.app.Logger.Info("All commands will be erased...")

	return helper.GetUserConfirmation("Continue to reset", i.app)
}
