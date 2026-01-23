package copyhandlers

import (
	"fmt"
	"runtime"

	appdata "github.com/fiwon123/crower/internal/data/app"
	executehandlers "github.com/fiwon123/crower/internal/handlers/execute"
)

// Copy file from a origin filepath to output folder
func CopyFile(filePath string, destFolder string, app *appdata.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("copy '%s' '%s'", filePath, destFolder), app)
	case "linux":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("\"cp '%s' '%s'\"", filePath, destFolder), app)
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	app.Logger.Info("output: ", "out", out)
	return nil
}

// Copy file from a origin folderpath to output folder
func CopyFolder(filePath string, destFolder string, app *appdata.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("xcopy '%s' '%s' /E /I", filePath, destFolder), app)
	case "linux":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("\"cp -r '%s' '%s'\"", filePath, destFolder), app)
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	app.Logger.Info("output: ", "out", out)
	return nil
}
