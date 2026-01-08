package handlers

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/pkg/utils"
)

// Delete command using key
func DeleteCommand(key string, app *app.Data) (*command.Data, bool) {
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
func DeleteFile(filePath string, app *app.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("del '%s'", filePath))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("\"rm '%s'\"", filePath))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	fmt.Println("result: ", out)
	return nil
}

// Delete folder from folderpath
func DeleteFolder(folderPath string, app *app.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("rmdir /s /q '%s'", folderPath))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf(`"rm -r '%s'"`, folderPath))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	fmt.Println("result: ", out)
	return nil
}

func DeleteHistoryContent(content history.Content, app *app.Data) (string, error) {
	ok := app.History.RemoveContent(content)
	if !ok {
		return "", fmt.Errorf("Content not found \n")
	}

	newDataPath := filepath.Join(app.HistoryFolderPath, content.File)
	err := utils.DeleteFile(newDataPath)
	if err != nil {
		return "", fmt.Errorf("Content not deleted %v \n", err)
	}

	return "Content deleted", nil
}
