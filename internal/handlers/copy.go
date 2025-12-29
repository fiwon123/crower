package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

func CopyFile(filePath string, destFolder string, app *app.Data) {
	switch runtime.GOOS {
	case "windows":
		PerformExecute(fmt.Sprintf("'copy %s %s'", filePath, destFolder))
	case "linux":
		PerformExecute(fmt.Sprintf("'cp %s %s'", filePath, destFolder))
	}

}

func CopyFolder(filePath string, destFolder string, app *app.Data) {
	switch runtime.GOOS {
	case "windows":
		PerformExecute(fmt.Sprintf("'xcopy %s %s /E /I'", filePath, destFolder))
	case "linux":
		PerformExecute(fmt.Sprintf("'cp -r %s %s'", filePath, destFolder))
	}

}
