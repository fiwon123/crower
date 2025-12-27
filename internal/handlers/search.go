package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

// Search based on user operational system(OS).
func Search(content string, app *app.Data) ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute(fmt.Sprintf(`start '' 'https://duckduckgo.com/?q=%s'`, content))
	case "linux":
		return PerformExecute(fmt.Sprintf(`'xdg-open "https://duckduckgo.com/?q=%s" '`, content))
	}

	return nil, nil
}

func SearchFile(currentPath string, content string, app *app.Data) ([]byte, error) {

	switch runtime.GOOS {
	case "windows":
		return PerformExecute(fmt.Sprintf(`dir "%s" /s /b /a-d | findstr "%s"`, currentPath, content))
	case "linux":
		return PerformExecute(fmt.Sprintf(`'find "%s" -type f -name "%s"'`, currentPath, content))
	}

	return nil, nil
}

func SearchFolder(currentPath string, content string, app *app.Data) ([]byte, error) {

	switch runtime.GOOS {
	case "windows":
		return PerformExecute(fmt.Sprintf(`dir "%s" /s /b /ad | findstr "%s"`, currentPath, content))
	case "linux":
		return PerformExecute(fmt.Sprintf(`'find "%s" -type d -name "%s"'`, currentPath, content))
	}

	return nil, nil
}

func SearchFileAndFolder(currentPath string, content string, app *app.Data) ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute(fmt.Sprintf(`dir "%s" /s /b | findstr "%s"`, currentPath, content))
	case "linux":
		return PerformExecute(fmt.Sprintf(`'find "%s" -name "%s"'`, currentPath, content))
	}

	return nil, nil

}
