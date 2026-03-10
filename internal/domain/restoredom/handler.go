package restoredom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/domain/createdom"
	"github.com/fiwon123/crower/internal/domain/updatedom"
)

type Handler struct {
	app           *app.Config
	handler       *updatedom.Handler
	createHandler *createdom.Handler
}

func NewHandler(app *app.Config) *Handler {

	createHandler := createdom.NewHandler(app)
	handler := updatedom.NewHandler(app)

	return &Handler{
		app:           app,
		handler:       handler,
		createHandler: createHandler,
	}
}

func (h *Handler) RestoreHistory(key string, content history.Content) (string, error) {
	command, err := h.app.FindCommand(key, content)
	if err != nil {
		return "", err
	}

	exists := h.app.AllCommandsByName.Exists(command.Name)
	// update if exists
	if exists {
		old, new, err := h.handler.UpdateCommand(key, "", nil, command.Exec)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("restored by updating: from %v to %v", old, new), nil

	}

	_, err = h.createHandler.CreateCommand(command.Name, command.AllAlias, command.Exec)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("restored by creating: %v", command), nil
}
