package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

func MoveFile(filePath string, destFolder string, app *app.Data) {
	switch runtime.GOOS {
	case "windows":
		PerformExecute(fmt.Sprintf("move '%s' '%s'", filePath, destFolder))
	case "linux":
		PerformExecute(fmt.Sprintf("mv '%s' '%s'", filePath, destFolder))
	}

}

func MoveFolder(folderPath string, destFolder string, app *app.Data) {
	switch runtime.GOOS {
	case "windows":
		PerformExecute(fmt.Sprintf("move '%s' '%s'", folderPath, destFolder))
	case "linux":
		PerformExecute(fmt.Sprintf("mv '%s' '%s'", folderPath, destFolder))
	}

}
