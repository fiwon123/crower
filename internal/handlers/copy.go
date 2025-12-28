package handlers

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

func CopyFile(currentFolder string, fileName string, destFolder string, app *app.Data) {
	fullPath := filepath.Join(currentFolder, fileName)
	switch runtime.GOOS {
	case "windows":
		PerformExecute(fmt.Sprintf("'copy %s %s'", fullPath, destFolder))
	case "linux":
		PerformExecute(fmt.Sprintf("'cp %s %s'", fullPath, destFolder))
	}

}

func CopyFolder(currentFolder string, folderName string, destFolder string, app *app.Data) {
	fullPath := filepath.Join(currentFolder, folderName)
	switch runtime.GOOS {
	case "windows":
		PerformExecute(fmt.Sprintf("'xcopy %s %s /E /I'", fullPath, destFolder))
	case "linux":
		PerformExecute(fmt.Sprintf("'cp -r %s %s'", fullPath, destFolder))
	}

}
