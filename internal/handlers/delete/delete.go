package deletehandlers

import (
	"fmt"
	"path/filepath"
	"runtime"

	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	historydata "github.com/fiwon123/crower/internal/data/history"

	executehandlers "github.com/fiwon123/crower/internal/handlers/execute"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

// Delete command using key
func DeleteCommand(key string, app *appdata.Data) (*commanddata.Data, bool) {
	command := app.AllCommandsByName.Get(key)
	if command == nil {
		command = app.AllCommandsByAlias.Get(key)
		if command == nil {
			return nil, false
		}
	}

	app.AllCommandsByName.Remove(command.Name)

	for _, alias := range command.AllAlias {
		app.AllCommandsByAlias.Remove(alias)
	}

	return command, true
}

// Delete file from filepath
func DeleteFile(filePath string, app *appdata.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("del '%s'", filePath), app)
	case "linux":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("\"rm '%s'\"", filePath), app)
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	app.Logger.Info("output: ", "out", out)
	return nil
}

// Delete folder from folderpath
func DeleteFolder(folderPath string, app *appdata.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("rmdir /s /q '%s'", folderPath), app)
	case "linux":
		out, err = executehandlers.PerformExecute(fmt.Sprintf(`"rm -r '%s'"`, folderPath), app)
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	app.Logger.Info("output: ", "out", out)
	return nil
}

func DeleteHistoryContent(content historydata.Content, app *appdata.Data) (string, error) {
	ok := app.History.RemoveContent(content)
	if !ok {
		return "", fmt.Errorf("Content not found \n")
	}

	newDataPath := filepath.Join(app.HistoryFolderPath, content.File)
	err := crowerutils.DeleteFile(newDataPath)
	if err != nil {
		return "", fmt.Errorf("Content not deleted %v \n", err)
	}

	return "Content deleted", nil
}
