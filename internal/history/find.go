package history

import (
	"fmt"
	"path/filepath"
	"slices"

	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	historydata "github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

func FindCommand(key string, content historydata.Content, app *appdata.Data) (*commanddata.Data, error) {
	allCommands := commanddata.NewMapData()
	newDataPath := filepath.Join(app.HistoryFolderPath, content.File)

	err := crowerutils.ReadToml(newDataPath, &allCommands)
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
