package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/pkg/utils"
)

// SearchBrowser based on user operational system(OS).
func SearchBrowser(content string, app *app.Data) (string, error) {

	switch runtime.GOOS {
	case "windows":
		if utils.IsURL(content) {
			return "", PerformExecuteStart(fmt.Sprintf(`start ' ' '%s' `, content), app)
		} else {
			return "", PerformExecuteStart(fmt.Sprintf(`start ' ' 'https://duckduckgo.com/?q=%s' `, content), app)
		}
	case "linux":
		if utils.IsURL(content) {
			return "", PerformExecuteStart(fmt.Sprintf(`xdg-open "%s" `, content), app)
		} else {
			return "", PerformExecuteStart(fmt.Sprintf(`xdg-open "https://duckduckgo.com/?q=%s" `, content), app)
		}

	}

	return "", nil
}

// Search files on folder path
func SearchFile(currentPath string, content string, app *app.Data) (string, error) {

	switch runtime.GOOS {
	case "windows":
		return PerformExecute(fmt.Sprintf(`dir "%s" /s /b /a-d | findstr "%s"`, currentPath, content), app)
	case "linux":
		return PerformExecute(fmt.Sprintf(`find "%s" -type f -name "%s"`, currentPath, content), app)
	}

	return "", nil
}

// Search folders on folder path
func SearchFolder(currentPath string, content string, app *app.Data) (string, error) {

	switch runtime.GOOS {
	case "windows":
		return PerformExecute(fmt.Sprintf(`dir "%s" /s /b /ad | findstr "%s"`, currentPath, content), app)
	case "linux":
		return PerformExecute(fmt.Sprintf(`find "%s" -type d -name "%s"`, currentPath, content), app)
	}

	return "", nil
}

// Search files and folder on folder path
func SearchFileAndFolder(currentPath string, content string, app *app.Data) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute(fmt.Sprintf(`dir "%s" /s /b | findstr "%s"`, currentPath, content), app)
	case "linux":
		return PerformExecute(fmt.Sprintf(`find "%s" -name "%s"`, currentPath, content), app)
	}

	return "", nil

}
