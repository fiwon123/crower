package history

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/pkg/utils"
)

func FindCommand(key string, content history.Content, app *app.Data) (*command.Data, error) {
	allCommands := command.NewMapData()
	newDataPath := filepath.Join(app.HistoryFolderPath, content.File)

	err := utils.ReadToml(newDataPath, &allCommands)
	if err != nil {
		return nil, err
	}

	for _, command := range allCommands {
		if command.Name == key {
			return &command, nil
		} else {
			if slices.Contains(command.AllAlias, key) {
				return &command, nil
			}
		}
	}

	return nil, fmt.Errorf("Not Found")
}
