package movehandlers

import (
	"fmt"
	"runtime"

	appdata "github.com/fiwon123/crower/internal/data/app"
	executehandlers "github.com/fiwon123/crower/internal/handlers/execute"
)

// Move file from origin path to output folder path
func MoveFile(filePath string, destFolder string, app *appdata.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("move '%s' '%s'", filePath, destFolder), app)
	case "linux":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("\"mv '%s' '%s'\"", filePath, destFolder), app)
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	app.Logger.Info("output: ", "out", out)
	return nil
}

// Move folder from origin path to output folder path
func MoveFolder(folderPath string, destFolder string, app *appdata.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("move '%s' '%s'", folderPath, destFolder), app)
	case "linux":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("\"mv '%s' '%s'\"", folderPath, destFolder), app)
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	app.Logger.Info("output: ", "out", out)
	return nil
}
