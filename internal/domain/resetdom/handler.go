package resetdom

import (
	"github.com/fiwon123/crower/internal/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
)

type Handler struct {
	app *app.Config
}

func NewHandler(app *app.Config) *Handler {

	return &Handler{
		app: app,
	}
}

// Reset all user cfg file.
func (h *Handler) Reset() {
	h.app.AllCommandsByName = commanddata.NewMapData()
}
