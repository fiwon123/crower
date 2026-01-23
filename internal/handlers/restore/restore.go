package restorehandlers

import (
	"fmt"

	appdata "github.com/fiwon123/crower/internal/data/app"
	dataHistory "github.com/fiwon123/crower/internal/data/history"
	createhandlers "github.com/fiwon123/crower/internal/handlers/create"
	updatehandlers "github.com/fiwon123/crower/internal/handlers/update"
	historyhelper "github.com/fiwon123/crower/internal/helper/history"
)

func RestoreHistory(key string, content dataHistory.Content, app *appdata.Data) (string, error) {
	command, err := historyhelper.FindCommand(key, content, app)
	if err != nil {
		return "", err
	}

	exists := app.AllCommandsByName.Exists(command.Name)
	// update if exists
	if exists {
		old, new, err := updatehandlers.UpdateCommand(key, "", nil, command.Exec, app)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("restored by updating: from %v to %v", old, new), nil

	}

	_, err = createhandlers.CreateCommand(command.Name, command.AllAlias, command.Exec, app)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("restored by creating: %v", command), nil
}
