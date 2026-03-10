package restoredom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/app"
	dataHistory "github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/domain/createdom"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
)

type Handler struct {
	app           *app.Config
	createHandler createdom.Handler
	updateHandler updatedom.Handler
}

func NewHandler(app *app.Config) *Handler {

	return &Handler{
		app: app,
	}
}

func (h *Handler) RestoreHistory(key string, content dataHistory.Content) (string, error) {
	command, err := historyhelper.FindCommand(key, content, h.app)
	if err != nil {
		return "", err
	}

	exists := h.app.AllCommandsByName.Exists(command.Name)
	// update if exists
	if exists {
		old, new, err := h.updatehandler.UpdateCommand(key, "", nil, command.Exec)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("restored by updating: from %v to %v", old, new), nil

	}

	_, err = h.createhandler.CreateCommand(command.Name, command.AllAlias, command.Exec)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("restored by creating: %v", command), nil
}
