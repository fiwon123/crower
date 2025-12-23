package handlers

import (
	"github.com/fiwon123/crower/internal/data"
)

func History(app *data.App) error {

	app.History.List()

	return nil
}
