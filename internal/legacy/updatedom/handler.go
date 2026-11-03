package updatedom

import (
	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/catalog"
)

// Handler applies an update to the loaded command set.
//
// The rules themselves live in the catalog aggregate; this type is the
// remaining bridge to the map-based configuration the operations still read. It
// disappears when `update` becomes a use case.
type Handler struct {
	app *app.Config
}

func NewHandler(app *app.Config) *Handler {

	return &Handler{
		app: app,
	}
}

// UpdateCommand changes the name, aliases or exec of a managed command.
//
// It returns the command before and after the change so the caller can record
// what happened in the history.
func (h *Handler) UpdateCommand(key string, newName string, newAlias []string, newExec string) (*command.Data, *command.Data, error) {

	result, err := h.app.Catalog().Update(key, catalog.Changes{
		Name:    newName,
		Aliases: newAlias,
		Exec:    newExec,
	})
	if err != nil {
		return nil, nil, err
	}

	h.app.ReplaceCatalog(result.Catalog)

	for _, warning := range result.Warnings {
		h.app.Logger.Warning(warning)
	}

	return toData(result.Old), toData(result.New), nil
}

// findCommand resolves a key the way every lookup does: name first, then
// alias. The input prompt and the handler share it, so the command shown for
// confirmation is the command that gets updated.
func findCommand(app *app.Config, key string) *command.Data {
	resolved, ok := app.Catalog().Resolve(key)
	if !ok {
		return nil
	}

	return toData(resolved)
}

// toData converts an aggregate command to the persisted shape.
func toData(c catalog.Command) *command.Data {
	return &command.Data{
		Name:     c.Name,
		AllAlias: c.Aliases,
		Exec:     c.Exec,
		Steps:    c.Steps,
	}
}
