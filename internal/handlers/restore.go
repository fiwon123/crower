package handlers

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	dataHistory "github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/history"
)

func RestoreHistory(key string, content dataHistory.Content, app *app.Data) (string, error) {
	command, err := history.FindCommand(key, content, app)
	if err != nil {
		return "", err
	}

	exists := app.AllCommandsByName.Exists(command.Name)
	// update if exists
	if exists {
		old, new, err := UpdateCommand(key, "", nil, command.Exec, app)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("restored by updating: from %v to %v", old, new), nil

	}

	_, err = CreateCommand(command.Name, command.AllAlias, command.Exec, app)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("restored by creating: %v", command), nil
}
