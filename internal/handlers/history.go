package handlers

import "github.com/fiwon123/crower/internal/data/app"

func History(app *app.Data) error {

	app.History.List()

	return nil
}
