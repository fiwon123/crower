package restoredom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/domain/catalog"
)

type Handler struct {
	app *app.Config
}

func NewHandler(app *app.Config) *Handler {
	return &Handler{app: app}
}

// RestoreHistory brings back the command a history entry holds: it updates the
// command when one with that name exists, and creates it otherwise.
func (h *Handler) RestoreHistory(key string, content history.Content) (string, error) {
	restored, err := h.app.FindCommand(key, content)
	if err != nil {
		return "", err
	}

	// update if exists
	if h.app.AllCommandsByName.Exists(restored.Name) {
		result, err := h.app.Catalog().Update(key, catalog.Changes{Exec: restored.Exec})
		if err != nil {
			return "", err
		}

		h.app.ReplaceCatalog(result.Catalog)

		for _, warning := range result.Warnings {
			h.app.Logger.Warning(warning)
		}

		return fmt.Sprintf("restored by updating: from %v to %v", toData(result.Old), toData(result.New)), nil
	}

	// The rules for adding a command -- empty fields, a taken name, a name that
	// is another command's alias -- live in the catalog, so restore uses the
	// same ones create does.
	commands, err := h.app.Catalog().Add(catalog.Command{
		Name:    restored.Name,
		Aliases: restored.AllAlias,
		Exec:    restored.Exec,
	})
	if err != nil {
		return "", err
	}

	h.app.ReplaceCatalog(commands)

	return fmt.Sprintf("restored by creating: %v", restored), nil
}

// toData converts an aggregate command to the shape this package prints.
func toData(c catalog.Command) *command.Data {
	return &command.Data{
		Name:     c.Name,
		AllAlias: c.Aliases,
		Exec:     c.Exec,
		Steps:    c.Steps,
	}
}
