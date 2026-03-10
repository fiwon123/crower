package searchdom

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/internal/domain/updatedom"
	"github.com/fiwon123/crower/pkg/utils"
)

type Handler struct {
	app            *app.Config
	handler        *updatedom.Handler
	executeHandler *executedom.Handler
}

func NewHandler(app *app.Config) *Handler {

	executeHandler := executedom.NewHandler(app)
	handler := updatedom.NewHandler(app)

	return &Handler{
		app:            app,
		handler:        handler,
		executeHandler: executeHandler,
	}
}

// SearchBrowser based on user operational system(OS).
func (h *Handler) SearchBrowser(content string) (string, error) {

	switch runtime.GOOS {
	case "windows":
		if utils.IsURL(content) {
			return "", h.executeHandler.PerformExecuteStart(fmt.Sprintf(`start ' ' '%s' `, content))
		} else {
			return "", h.executeHandler.PerformExecuteStart(fmt.Sprintf(`start ' ' 'https://duckduckgo.com/?q=%s' `, content))
		}
	case "linux":
		if utils.IsURL(content) {
			return "", h.executeHandler.PerformExecuteStart(fmt.Sprintf(`xdg-open "%s" `, content))
		} else {
			return "", h.executeHandler.PerformExecuteStart(fmt.Sprintf(`xdg-open "https://duckduckgo.com/?q=%s" `, content))
		}

	}

	return "", nil
}

// Search files on folder path
func (h *Handler) SearchFile(currentPath string, content string) (string, error) {

	switch runtime.GOOS {
	case "windows":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`dir "%s" /s /b /a-d | findstr "%s"`, currentPath, content))
	case "linux":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`find "%s" -type f -name "%s"`, currentPath, content))
	}

	return "", nil
}

// Search folders on folder path
func (h *Handler) SearchFolder(currentPath string, content string) (string, error) {

	switch runtime.GOOS {
	case "windows":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`dir "%s" /s /b /ad | findstr "%s"`, currentPath, content))
	case "linux":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`find "%s" -type d -name "%s"`, currentPath, content))
	}

	return "", nil
}

// Search files and folder on folder path
func (h *Handler) SearchFileAndFolder(currentPath string, content string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`dir "%s" /s /b | findstr "%s"`, currentPath, content))
	case "linux":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`find "%s" -name "%s"`, currentPath, content))
	}

	return "", nil

}
