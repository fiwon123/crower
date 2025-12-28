package handlers

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

func MoveFile(currentFolder string, fileName string, destFolder string, app *app.Data) {
	fullPath := filepath.Join(currentFolder, fileName)
	switch runtime.GOOS {
	case "windows":
		PerformExecute(fmt.Sprintf("'move %s %s'", fullPath, destFolder))
	case "linux":
		PerformExecute(fmt.Sprintf("'mv %s %s'", fullPath, destFolder))
	}

}

func MoveFolder(currentFolder string, folderName string, destFolder string, app *app.Data) {
	fullPath := filepath.Join(currentFolder, folderName)
	switch runtime.GOOS {
	case "windows":
		PerformExecute(fmt.Sprintf("'move %s %s'", fullPath, destFolder))
	case "linux":
		PerformExecute(fmt.Sprintf("'mv %s %s'", fullPath, destFolder))
	}

}
